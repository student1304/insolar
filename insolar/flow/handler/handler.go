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
	"errors"
	"fmt"
	"strconv"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/flow/internal/pulse"
	"github.com/insolar/insolar/insolar/flow/internal/thread"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

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
		ReplyTo: make(chan bus.Reply, 1),
		Parcel:  parcel,
	}
	defer func() {
		close(msg.ReplyTo)
	}()

	ctx = pulse.ContextWith(ctx, parcel.Pulse())

	f := thread.NewThread(msg, h.controller)
	err := f.Run(ctx, h.handles.present(msg))

	var rep bus.Reply
	select {
	case rep = <-msg.ReplyTo:
		return rep.Reply, rep.Err
	default:
	}

	if err != nil {
		return nil, err
	}

	return nil, errors.New("no reply from handler")
}

func (h *Handler) InnerSubscriber(watermillMsg *message.Message) ([]*message.Message, error) {
	msg := bus.Message{
		WatermillMsg: watermillMsg,
	}

	ctx := watermillMsg.Context()
	logger := inslogger.FromContext(ctx)
	go func() {
		f := thread.NewThread(msg, h.controller)
		err := f.Run(ctx, h.handles.present(msg))
		if err != nil {
			logger.Error("Handling failed", err)
		}
	}()
	return nil, nil
}

func (h *Handler) Process(ctx context.Context, msg *message.Message, pub message.Publisher) error {
	msgBus := bus.Message{
		Publisher:    pub,
		WatermillMsg: msg,
		ReplyTo:      make(chan bus.Reply),
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
		resAsMsg.Metadata.Set("Type", "Reply")
		fmt.Println("get Reply with UUid:", resAsMsg.UUID)
		err = pub.Publish("outbound", resAsMsg)
		if err != nil {
			fmt.Println("All was bad, really bad", err)
		}
	}(msgBus, pub)
	return nil
}
