//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package common

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"runtime/debug"
	"sync"

	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/logicrunner"
	"github.com/insolar/insolar/logicrunner/artifacts"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type RPCMethods struct {
	lr *logicrunner.LogicRunner
}

func recoverRPC(err *error) {
	if r := recover(); r != nil {
		// Global logger is used because there is no access to context here
		log.Errorf("Recovered panic:\n%s", string(debug.Stack()))
		if err != nil {
			if *err == nil {
				*err = errors.New(fmt.Sprint(r))
			} else {
				*err = errors.New(fmt.Sprint(*err, r))
			}
		}
	}
}

// GetCode is an RPC retrieving a code by its reference
func (gpr *RPCMethods) GetCode(req rpctypes.UpGetCodeReq, reply *rpctypes.UpGetCodeResp) (err error) {
	defer recoverRPC(&err)
	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context
	inslogger.FromContext(ctx).Debug("In RPC.GetCode ....")

	am := gpr.lr.ArtifactManager

	ctx, span := instracer.StartSpan(ctx, "service.GetCode")
	defer span.End()

	codeDescriptor, err := am.GetCode(ctx, req.Code)
	if err != nil {
		return err
	}
	reply.Code, err = codeDescriptor.Code()
	if err != nil {
		return err
	}
	return nil
}

// RouteCall routes call from a contract to a contract through event bus.
func (gpr *RPCMethods) RouteCall(req rpctypes.UpRouteReq, rep *rpctypes.UpRouteResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)

	if os.ExecutionState.Current.LogicContext.Immutable {
		return errors.New("Try to call route from immutable method")
	}

	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	// TODO: delegation token

	es.nonce++

	msg := &message.CallMethod{
		Request: record.Request{
			Caller:          req.Callee,
			CallerPrototype: req.CalleePrototype,
			Nonce:           es.nonce,

			Immutable: req.Immutable,

			Object:    &req.Object,
			Prototype: &req.Prototype,
			Method:    req.Method,
			Arguments: req.Arguments,
		},
	}

	if !req.Wait {
		msg.ReturnMode = record.ReturnNoWait
	}

	res, err := gpr.lr.ContractRequester.CallMethod(ctx, msg)
	if err != nil {
		return err
	}

	if req.Wait {
		rep.Result = res.(*reply.CallMethod).Result
	}

	return nil
}

// SaveAsChild is an RPC saving data as memory of a contract as child a parent
func (gpr *RPCMethods) SaveAsChild(req rpctypes.UpSaveAsChildReq, rep *rpctypes.UpSaveAsChildResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	es.nonce++

	msg := &message.CallMethod{
		Request: record.Request{
			Caller:          req.Callee,
			CallerPrototype: req.CalleePrototype,
			Nonce:           es.nonce,

			CallType:  record.CTSaveAsChild,
			Base:      &req.Parent,
			Prototype: &req.Prototype,
			Method:    req.ConstructorName,
			Arguments: req.ArgsSerialized,
		},
	}

	ref, err := gpr.lr.ContractRequester.CallConstructor(ctx, msg)

	rep.Reference = ref

	return err
}

// SaveAsDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPCMethods) SaveAsDelegate(req rpctypes.UpSaveAsDelegateReq, rep *rpctypes.UpSaveAsDelegateResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	es.nonce++

	msg := &message.CallMethod{
		Request: record.Request{
			Caller:          req.Callee,
			CallerPrototype: req.CalleePrototype,
			Nonce:           es.nonce,

			CallType:  record.CTSaveAsDelegate,
			Base:      &req.Into,
			Prototype: &req.Prototype,
			Method:    req.ConstructorName,
			Arguments: req.ArgsSerialized,
		},
	}

	ref, err := gpr.lr.ContractRequester.CallConstructor(ctx, msg)

	rep.Reference = ref
	return err
}

var iteratorMap = make(map[string]artifacts.RefIterator)
var iteratorMapLock = sync.RWMutex{}
var iteratorBuffSize = 1000

// GetObjChildrenIterator is an RPC returns an iterator over object children with specified prototype
func (gpr *RPCMethods) GetObjChildrenIterator(
	req rpctypes.UpGetObjChildrenIteratorReq,
	rep *rpctypes.UpGetObjChildrenIteratorResp,
) (
	err error,
) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	am := gpr.lr.ArtifactManager
	iteratorID := req.IteratorID

	iteratorMapLock.RLock()
	iterator, ok := iteratorMap[iteratorID]
	iteratorMapLock.RUnlock()

	if !ok {
		newIterator, err := am.GetChildren(ctx, req.Object, nil)
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get children")
		}

		id, err := uuid.NewV4()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't generate UUID")
		}

		iteratorID = id.String()

		iteratorMapLock.Lock()
		iterator, ok = iteratorMap[iteratorID]
		if !ok {
			iteratorMap[iteratorID] = newIterator
			iterator = newIterator
		}
		iteratorMapLock.Unlock()
	}

	iter := iterator

	rep.Iterator.ID = iteratorID
	rep.Iterator.CanFetch = iter.HasNext()
	for len(rep.Iterator.Buff) < iteratorBuffSize && iter.HasNext() {
		r, err := iter.Next()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get Next")
		}
		rep.Iterator.CanFetch = iter.HasNext()

		o, err := am.GetObject(ctx, *r)

		if err != nil {
			if err == insolar.ErrDeactivated {
				continue
			}
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't call GetObject on Next")
		}
		protoRef, err := o.Prototype()
		if err != nil {
			return errors.Wrap(err, "[ GetObjChildrenIterator ] Can't get prototype reference")
		}

		if protoRef.Equal(req.Prototype) {
			rep.Iterator.Buff = append(rep.Iterator.Buff, *r)
		}
	}

	if !iter.HasNext() {
		iteratorMapLock.Lock()
		delete(iteratorMap, rep.Iterator.ID)
		iteratorMapLock.Unlock()
	}

	return nil
}

// GetDelegate is an RPC saving data as memory of a contract as child a parent
func (gpr *RPCMethods) GetDelegate(req rpctypes.UpGetDelegateReq, rep *rpctypes.UpGetDelegateResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	ctx := es.Current.Context

	am := gpr.lr.ArtifactManager
	ref, err := am.GetDelegate(ctx, req.Object, req.OfType)
	if err != nil {
		return err
	}
	rep.Object = *ref
	return nil
}

// DeactivateObject is an RPC saving data as memory of a contract as child a parent
func (gpr *RPCMethods) DeactivateObject(req rpctypes.UpDeactivateObjectReq, rep *rpctypes.UpDeactivateObjectResp) (err error) {
	defer recoverRPC(&err)

	os := gpr.lr.MustObjectState(req.Callee)
	es := os.MustModeState(req.Mode)
	es.deactivate = true
	return nil
}

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	server    *rpc.Server
	methods   *RPCMethods
	listener  net.Listener
	proto     string
	listen    string
	isStarted bool
}

func NewRPC(ctx context.Context, lr *logicrunner.LogicRunner) *RPC {
	rpcService := &RPC{
		server:  rpc.NewServer(),
		methods: &RPCMethods{},
		proto:   lr.Cfg.RPCListen,
		listen:  lr.Cfg.RPCProtocol,
	}
	if err := rpcService.server.Register(rpcService.methods); err != nil {
		panic("Fail to register LogicRunner RPC Service: " + err.Error())
	}

	return rpcService
}

// StartRPC starts RPC server for isolated executors to use
func (rpc *RPC) Start(ctx context.Context) {
	var err error
	logger := inslogger.FromContext(ctx)

	rpc.listener, err = net.Listen(rpc.proto, rpc.listen)
	if err != nil {
		logger.Fatalf("couldn't setup listener on %q over %q: %s", rpc.listen, rpc.proto, err)
	}

	logger.Infof("starting LogicRunner RPC service on %q over %q", rpc.listen, rpc.proto)
	rpc.isStarted = true

	go func() {
		rpc.server.Accept(rpc.listener)
		logger.Info("LogicRunner RPC service stopped")
	}()
}

func (rpc *RPC) Stop(_ context.Context) error {
	if rpc.isStarted {
		rpc.isStarted = false
		if err := rpc.listener.Close(); err != nil {
			return err
		}
	}
	return nil
}
