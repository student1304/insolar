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

package handler

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/pkg/errors"
)

const handleTimeout = 10 * time.Second

type Handler struct {
	handles struct {
		present flow.MakeHandle
	}
	controller *thread.Controller
}

func NewHandler(present flow.MakeHandle) *Handler {
	h := &Handler{
		controller: thread.NewController(),
	}
	h.handles.present = present
	return h
}

// ChangePulse is a handle for pulse change vent.
func (h *Handler) ChangePulse(ctx context.Context, pulse insolar.Pulse) {
	h.controller.Pulse()
}

func (h *Handler) WrapBusHandle(ctx context.Context, parcel insolar.Parcel) (insolar.Reply, error) {
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply),
		Parcel:  parcel,
	}
	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", parcel.Pulse()))
	ctx = pulse.ContextWith(ctx, parcel.Pulse())
	go func() {
		f := thread.NewThread(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		if err != nil {
			select {
			case msg.ReplyTo <- bus.Reply{Err: err}:
			default:
			}
			logger.Error("Handling failed", err)
		}
	}()
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	case <-time.After(handleTimeout):
		return nil, errors.New("handler timeout")
	}
}

func (h *Handler) Innner(msgWM *message.Message) ([]*message.Message, error) {
	msg := bus.Message{
		ReplyTo: make(chan bus.Reply),
		Msg:     msgWM,
	}
	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", parcel.Pulse()))
	ctx = pulse.ContextWith(ctx, parcel.Pulse())
	go func() {
		f := thread.NewThread(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		if err != nil {
			select {
			case msg.ReplyTo <- bus.Reply{Err: err}:
			default:
			}
			logger.Error("Handling failed", err)
		}
	}()
	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	case <-time.After(handleTimeout):
		return nil, errors.New("handler timeout")
	}
	return nil, nil
}

func (h *Handler) Process(ctx context.Context, msg *message.Message, pub message.Publisher) error {
	msgBus := bus.Message{
		Publisher: pub,
		Msg:       msg,
		ReplyTo:   make(chan bus.Reply),
	}
	// fix it
	pStr := msg.Metadata.Get("Pulse")
	fmt.Println("pInt is ", pStr)
	u64, err := strconv.ParseUint(pStr, 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	pInt := uint32(u64)
	p := insolar.PulseNumber(pInt)

	ctx, logger := inslogger.WithField(ctx, "pulse", fmt.Sprintf("%d", p))
	ctx = pulse.ContextWith(ctx, p)
	go func() {
		f := thread.NewThread(msgBus, h.controller)
		err := f.Run(ctx, h.handles.present(msgBus))
		if err != nil {
			logger.Error("Handling failed", err)
		}
	}()
	go func(msg bus.Message, pub message.Publisher) {
		rep := <-msg.ReplyTo
		rd, err := reply.Serialize(rep.Reply)
		fmt.Println("get Reply:", rep.Reply.Type())
		if err != nil {
			fmt.Println("All was bad, really bad", err)
		}
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(rd)
		if err != nil {
			fmt.Println("All was bad, really bad", err)
		}
		resInBytes := buf.Bytes()
		resAsMsg := message.NewMessage(watermill.NewUUID(), resInBytes)
		fmt.Println("get Reply with UUid:", resAsMsg.UUID)
		err = pub.Publish("outbound", resAsMsg)
		if err != nil {
			fmt.Println("All was bad, really bad", err)
		}
	}(msgBus, pub)
	return nil
}
