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

package dispatcher

import (
	"bytes"
	"context"
	"errors"
	"runtime"
	"strconv"
	"sync/atomic"

	"github.com/insolar/insolar/log"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/instrumentation/inslogger"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
)

const TraceIDField = "TraceID"

type Dispatcher struct {
	handles struct {
		present flow.MakeHandle
		future  flow.MakeHandle
	}
	controller         *thread.Controller
	currentPulseNumber uint32
}

func NewDispatcher(present flow.MakeHandle, future flow.MakeHandle) *Dispatcher {
	log.Debug("NEW DISPATCHER CREATED") // TODO FIXME Dispatcher is created in 4 different places!
	d := &Dispatcher{
		controller: thread.NewController(),
	}
	d.handles.present = present
	d.handles.future = future
	d.currentPulseNumber = insolar.FirstPulseNumber
	return d
}

func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

// ChangePulse is a handle for pulse change vent.
func (d *Dispatcher) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	// panic("ChangePulse")
	/*
	   github.com/insolar/insolar/insolar/flow/dispatcher.(*Dispatcher).ChangePulse(0xc0001f93c0, 0x1c5dfc0, 0xc0003c2060, 0x10002, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	   	/Users/eax/go/src/github.com/insolar/insolar/insolar/flow/dispatcher/dispatcher.go:72 +0x39
	   github.com/insolar/insolar/logicrunner.(*LogicRunner).OnPulse(0xc0005f6000, 0x1c5dfc0, 0xc0003c2060, 0x10002, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner.go:847 +0xbe
	   github.com/insolar/insolar/logicrunner/pulsemanager.(*PulseManager).Set(0xc0005f6280, 0x1c5dfc0, 0xc0003c2060, 0x10002, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, ...)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/pulsemanager/pulsemanager.go:97 +0x37b
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).incrementPulseHelper(0xc0000ff180, 0x1c5df40, 0xc0000bc010, 0x1c5e5c0, 0xc0005f6000, 0x1c54280, 0xc0005f6280)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:194 +0x1ba
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).PrepareLrAmCbPm(0xc0000ff180, 0x1, 0x1, 0x10046f7, 0x571f23a686c38c4d, 0x1a44c60, 0xc000164228, 0xc000077858, 0x100e078, 0xc0001641a0)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:177 +0xe00
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).TestBasicNotificationCallError(0xc0000ff180)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:650 +0xb4
	   reflect.Value.call(0xc0000402a0, 0xc0000ca608, 0x13, 0x1b3337b, 0x4, 0xc000077f80, 0x1, 0x1, 0xc000257eb8, 0x78, ...)
	*/
	log.Debug("[GID ", getGID(), "] [SELF ", d, "] WrapBusHandle-CHANGE PULSE ", uint32(pulse.PulseNumber))
	d.controller.Pulse()
	atomic.StoreUint32(&d.currentPulseNumber, uint32(pulse.PulseNumber)) // TODO FIXME WTF???
	log.Debug("[GID ", getGID(), "] [SELF ", d, "NEW PULSE = ", atomic.LoadUint32(&d.currentPulseNumber))
}

func (d *Dispatcher) getHandleByPulse(msgPulseNumber insolar.PulseNumber) flow.MakeHandle {
	// panic("GetHandleByPulse")
	/*
	   	/usr/local/go/src/runtime/panic.go:513 +0x1b9
	   github.com/insolar/insolar/insolar/flow/dispatcher.(*Dispatcher).getHandleByPulse(0xc00033dcc0, 0xc000010002, 0xc000010002)
	   	/Users/eax/go/src/github.com/insolar/insolar/insolar/flow/dispatcher/dispatcher.go:95 +0x39
	   github.com/insolar/insolar/insolar/flow/dispatcher.(*Dispatcher).WrapBusHandle(0xc00033dcc0, 0x1c5df20, 0xc000581290, 0x1c64720, 0xc0000d8840, 0xc0001852d0, 0xc000078310, 0x1, 0x1)
	   	/Users/eax/go/src/github.com/insolar/insolar/insolar/flow/dispatcher/dispatcher.go:113 +0x1c0
	   github.com/insolar/insolar/insolar/flow/dispatcher.(*Dispatcher).WrapBusHandle-fm(0x1c5df20, 0xc000581260, 0x1c64720, 0xc0000d8840, 0xc000581260, 0x18a7fdfc86b034da, 0x1c094919a0cd2288, 0xdfb4617642c987bc)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner.go:280 +0x52
	   github.com/insolar/insolar/testutils/testmessagebus.(*TestMessageBus).Send(0xc0005748a0, 0x1c5df20, 0xc000581260, 0x1c60360, 0xc0005b0600, 0x0, 0x0, 0x0, 0x0, 0x0)
	   	/Users/eax/go/src/github.com/insolar/insolar/testutils/testmessagebus/testmessagebus.go:166 +0x863
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).incrementPulseHelper(0xc0001ccaa0, 0x1c5dea0, 0xc0000bc010, 0x1c5e520, 0xc00018e000, 0x1c541e0, 0xc00018e280)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:203 +0x4a0
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).PrepareLrAmCbPm(0xc0001ccaa0, 0x1, 0x1, 0x10046f7, 0xe2aeb1beb1cb757c, 0x1a44ba0, 0xc000182228, 0xc000072858, 0x100e078, 0xc0001821a0)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:177 +0xe00
	   github.com/insolar/insolar/logicrunner.(*LogicRunnerFuncSuite).TestBasicNotificationCallError(0xc0001ccaa0)
	   	/Users/eax/go/src/github.com/insolar/insolar/logicrunner/logicrunner_test.go:650 +0xb4
	*/
	log.Debug("[GID ", getGID(), "] [SELF", d, " WrapBusHandle-getHandleByPulse msgPulseNumber = ", msgPulseNumber, "current = ", atomic.LoadUint32(&d.currentPulseNumber))
	if uint32(msgPulseNumber) > atomic.LoadUint32(&d.currentPulseNumber) {
		return d.handles.future
	}
	return d.handles.present
}

func (d *Dispatcher) WrapBusHandle(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	log.Debugf("WrapBusHandle-BEGIN, type = %v, pulse = %v", parcel.Type(), parcel.Pulse())
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply, 1),
		Parcel:  parcel,
	}

	ctx = pulse.ContextWith(ctx, parcel.Pulse())

	f := thread.NewThread(msg, d.controller)
	handle := d.getHandleByPulse(parcel.Pulse())

	log.Debugf("WrapBusHandle-BEFORE-RUN, handle = %v", handle)

	err := f.Run(ctx, handle(msg))
	log.Debug("WrapBusHandle-BEFORE-SELECT")
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	default:
	}

	log.Debug("WrapBusHandle-AFTER-SELECT")

	if err != nil {
		return nil, err
	}

	log.Debug("WrapBusHandle-END")
	return nil, errors.New("no reply from handler")
}

func (d *Dispatcher) InnerSubscriber(watermillMsg *message.Message) ([]*message.Message, error) {
	msg := bus.Message{
		WatermillMsg: watermillMsg,
	}

	ctx := context.Background()
	ctx = inslogger.ContextWithTrace(ctx, watermillMsg.Metadata.Get(TraceIDField))
	logger := inslogger.FromContext(ctx)
	go func() {
		f := thread.NewThread(msg, d.controller)
		err := f.Run(ctx, d.handles.present(msg))
		if err != nil {
			logger.Error("Handling failed", err)
		}
	}()
	return nil, nil
}

// Process handles incoming message.
func (d *Dispatcher) Process(msg *message.Message) ([]*message.Message, error) {
	return nil, nil
}
