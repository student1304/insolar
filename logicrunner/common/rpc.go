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
	"net"
	"net/rpc"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
)

type LogicRunnerRPCStub interface {
	GetCode(rpctypes.UpGetCodeReq, *rpctypes.UpGetCodeResp) error
	RouteCall(rpctypes.UpRouteReq, *rpctypes.UpRouteResp) error
	SaveAsChild(rpctypes.UpSaveAsChildReq, *rpctypes.UpSaveAsChildResp) error
	SaveAsDelegate(rpctypes.UpSaveAsDelegateReq, *rpctypes.UpSaveAsDelegateResp) error
	GetObjChildrenIterator(rpctypes.UpGetObjChildrenIteratorReq, *rpctypes.UpGetObjChildrenIteratorResp) error
	GetDelegate(rpctypes.UpGetDelegateReq, *rpctypes.UpGetDelegateResp) error
	DeactivateObject(rpctypes.UpDeactivateObjectReq, *rpctypes.UpDeactivateObjectResp) error
}

// func recoverRPC(err *error) {
// 	if r := recover(); r != nil {
// 		// Global logger is used because there is no access to context here
// 		log.Errorf("Recovered panic:\n%s", string(debug.Stack()))
// 		if err != nil {
// 			if *err == nil {
// 				*err = errors.New(fmt.Sprint(r))
// 			} else {
// 				*err = errors.New(fmt.Sprint(*err, r))
// 			}
// 		}
// 	}
// }

// RPC is a RPC interface for runner to use for various tasks, e.g. code fetching
type RPC struct {
	server    *rpc.Server
	methods   LogicRunnerRPCStub
	listener  net.Listener
	proto     string
	listen    string
	isStarted bool
}

func NewRPC(_ context.Context, lr LogicRunnerRPCStub, cfg *configuration.LogicRunner) *RPC {
	rpcService := &RPC{
		server:  rpc.NewServer(),
		methods: lr,
		proto:   cfg.RPCProtocol,
		listen:  cfg.RPCListen,
	}
	if err := rpcService.server.RegisterName("RPC", rpcService.methods); err != nil {
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
