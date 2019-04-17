/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package bus

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/infrastructure/gochannel"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow/handler"
	insMsg "github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/ledger/storage/pulse"
)

const deliverIncomeMsg = "Bus.Income"

type Bus struct {
	Network        insolar.Network        `inject:""`
	JetCoordinator insolar.JetCoordinator `inject:""`
	PulseAccessor  pulse.Accessor         `inject:""`
	ParcelFactory  insMsg.ParcelFactory   `inject:""`
	NodeNetwork    insolar.NodeNetwork    `inject:""`

	Handler *handler.Handler

	Pub         message.Publisher
	sub         message.Subscriber
	resultMap   map[string]*future
	resultMutex sync.RWMutex
}

func NewBus(
	handler *handler.Handler,
) *Bus {
	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)
	// logger := watermill.NewStdLogger(false, false)
	// router, err := message.NewRouter(message.RouterConfig{}, logger)
	// if err != nil {
	// 	panic(err)
	// }
	// router.AddPlugin(plugin.SignalsHandler)
	// router.AddMiddleware(
	// 	// correlation ID will copy correlation id from consumed message metadata to produced messages
	// 	middleware.CorrelationID,
	// )
	bus := &Bus{
		Pub:       pubSub,
		sub:       pubSub,
		Handler:   handler,
		resultMap: make(map[string]*future),
	}
	// router.AddNoPublisherHandler(
	// 	"inbound",                // handler name, must be unique
	// 	"example.topiinboundc_1", // topic from which we will read events
	// 	pubSub,
	// 	nil,
	// )
	// router.AddHandler(
	// 	"struct_handler", // handler name, must be unique
	// 	"outbound",       // topic from which we will read events
	// 	pubSub,
	// 	"inbound", // topic to which we will publish event
	// 	pubSub,
	// 	nil,
	// )
	return bus
}

// Start initializes message bus.
func (b *Bus) Start(ctx context.Context) error {
	fmt.Println("bus was started")
	// b.Network.RemoteProcedureRegister(deliverWatermillMsg, b.get)
	b.Network.RemoteProcedureRegister(deliverIncomeMsg, b.toIncome)

	inMessages, err := b.sub.Subscribe(context.Background(), "inbound")
	if err != nil {
		panic(err)
	}

	outMessages, err := b.sub.Subscribe(context.Background(), "outbound")
	if err != nil {
		panic(err)
	}

	go b.processIncome(ctx, inMessages)
	go b.processOutcome(ctx, outMessages)

	return nil
}

func (b *Bus) toIncome(ctx context.Context, args [][]byte) ([]byte, error) {
	inslogger.FromContext(ctx).Debug("MessageBus.toIncome starts ...")
	fmt.Println("deliverIncomeMsg income")
	if len(args) < 1 {
		return nil, errors.New("need exactly one argument when mb.deliver()")
	}
	parcel, err := insMsg.DeserializeParcel(bytes.NewBuffer(args[0]))
	if err != nil {
		return nil, err
	}

	// parcelCtx := parcel.Context(context.Background()) // use ctx when network provide context
	// inslogger.FromContext(ctx).Debugf("MessageBus.deliver after deserialize msg. Msg Type: %s", parcel.Type())

	// mb.globalLock.RLock()
	//
	// if err = mb.checkPulse(parcelCtx, parcel, true); err != nil {
	// 	mb.globalLock.RUnlock()
	// 	return nil, err
	// }
	//
	// if err = mb.checkParcel(parcelCtx, parcel); err != nil {
	// 	mb.globalLock.RUnlock()
	// 	return nil, err
	// }
	// mb.globalLock.RUnlock()

	// create msg from network bytes
	wmMsg := parcel.Message().(*insolar.Watermill)
	msg := wmMsg.Msg
	msg.Metadata.Set("Sender", parcel.GetSender().String())
	fmt.Println("hi love, get type:", msg.Metadata.Get("Type"))
	fmt.Println("hi love, get mets:", msg.Metadata, msg.UUID)
	fmt.Println("hi love, get pays:", msg.Payload)
	err = b.Pub.Publish("inbound", &msg)
	if err != nil {
		// 	do staff
	}

	rd, err := reply.Serialize(&reply.OK{})
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(rd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *Bus) processIncome(ctx context.Context, messages <-chan *message.Message) {
	for msg := range messages {
		log.Printf("received message: %s, payload: %s", msg.UUID, string(msg.Payload))

		b.resultMutex.RLock()
		f, ok := b.resultMap[msg.UUID]
		b.resultMutex.RUnlock()
		if !ok {
			fmt.Println(b.Handler)
			_ = b.Handler.Process(ctx, msg, b.Pub)
			msg.Ack()
			continue
		}

		fmt.Println("msg.UUID love,", msg.UUID)
		res, err := reply.Deserialize(bytes.NewBuffer(msg.Payload))
		if err != nil {
			fmt.Println("lol kek sorry", err)
			// 	do staff
		}
		f.SetResult(res)

		// we need to Acknowledge that we received and processed the message,
		// otherwise we will not receive next message
		msg.Ack()
	}
}

func (b *Bus) processOutcome(ctx context.Context, messages <-chan *message.Message) {
	for msg := range messages {

		sender := msg.Metadata.Get("Sender")
		if sender != "" {
			serviceData := insMsg.ServiceData{
				LogTraceID:    inslogger.TraceID(ctx),
				LogLevel:      inslogger.GetLoggerLevel(ctx),
				TraceSpanData: instracer.MustSerialize(ctx),
			}

			parcelForNet := &insMsg.Parcel{
				Msg: &insolar.Watermill{
					Msg: *msg,
				},
				// Signature:   signature.Bytes(),
				Sender:      sender,
				Token:       parcel.DelegationToken(),
				PulseNumber: parcel.Pulse(),
				ServiceData: serviceData,
			}
			fmt.Println("deliverIncomeMsg was send")
			_, err = b.Network.SendMessage(nodes[0], deliverIncomeMsg, parcelForNet)
		}

		parcel, err := insMsg.DeserializeParcel(bytes.NewBuffer(msg.Payload))
		if err != nil {
			fmt.Println("lol kek err,", err)
			// 	do staff
		}

		var (
			nodes []insolar.Reference
		)
		if msg.Metadata.Get("Receiver") != "" {
			ref, err := insolar.NewReferenceFromBase58(msg.Metadata.Get("Receiver"))
			if err != nil {
				// 	do staff
			}
			nodes = []insolar.Reference{*ref}
			fmt.Println("kek no rec", nodes)
		} else {
			// TODO: send to all actors of the role if nil Target
			fmt.Println("lol kek parcel - ", parcel)
			fmt.Println("lol kek", parcel.Type().String())
			target := parcel.DefaultTarget()
			// FIXME: @andreyromancev. 21.12.18. Temp hack. All messages should have a default target.
			if target == nil {
				target = &insolar.Reference{}
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

			nodes, err = b.JetCoordinator.QueryRole(ctx, parcel.DefaultRole(), *target.Record(), p)
			if err != nil {
				// 	do staff
				fmt.Println("kek err", err)
			}
			fmt.Println("kek yes rec", nodes)
		}
		// f := newFuture()
		// payload := insMsg.ParcelToBytes(parcel)
		// msg := message.NewMessage(watermill.NewUUID(), payload)
		// b.resultMutex.Lock()
		// b.resultMap[msg.UUID] = f
		// b.resultMutex.Unlock()
		// b.Network.SendMessage()

		// if len(nodes) > 1 {
		// 	cascade := insolar.Cascade{
		// 		NodeIds:           nodes,
		// 		Entropy:           currentPulse.Entropy,
		// 		ReplicationFactor: 2,
		// 	}
		// 	err := mb.Network.SendCascadeMessage(cascade, deliverRPCMethodName, parcelForNet)
		// 	return nil, err
		// }

		// Short path when sending to self node. Skip serialization
		origin := b.NodeNetwork.GetOrigin()
		fmt.Println("Hi love, origin:", origin, ", nodes[0]:", nodes[0])
		if nodes[0].Equal(origin.ID()) {
			err = b.Pub.Publish("inbound", msg)
			if err != nil {
				// 	do staff
			}
			msg.Ack()
			continue
		}

		// lol kek fix it
		serviceData := insMsg.ServiceData{
			LogTraceID:    inslogger.TraceID(ctx),
			LogLevel:      inslogger.GetLoggerLevel(ctx),
			TraceSpanData: instracer.MustSerialize(ctx),
		}

		parcelForNet := &insMsg.Parcel{
			Msg: &insolar.Watermill{
				Msg: *msg,
			},
			// Signature:   signature.Bytes(),
			Sender:      parcel.GetSender(),
			Token:       parcel.DelegationToken(),
			PulseNumber: parcel.Pulse(),
			ServiceData: serviceData,
		}
		fmt.Println("deliverIncomeMsg was send")
		_, err = b.Network.SendMessage(nodes[0], deliverIncomeMsg, parcelForNet)
		if err != nil {
			fmt.Println("err in deliverIncomeMsg:", err)
			// 	do staff
		}

		// we need to Acknowledge that we received and processed the message,
		// otherwise we will not receive next message
		msg.Ack()
	}
}

func (b *Bus) Send(ctx context.Context, msg insolar.Message, ops *insolar.MessageSendOptions) (insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "MessageBus.Send "+msg.Type().String())
	defer span.End()
	wait, err := b.SendAsync(ctx, msg, ops)
	if err != nil {
		return nil, err
	}
	select {
	case res := <-wait:
		return res, nil
	case <-time.After(time.Second * 3):
		return nil, errors.New("timeout from async send")
	}
}

// Send an `Message` and get a `Value` or error from remote host.
func (b *Bus) SendAsync(ctx context.Context, msg insolar.Message, ops *insolar.MessageSendOptions) (<-chan insolar.Reply, error) {
	ctx, span := instracer.StartSpan(ctx, "MessageBus.Send "+msg.Type().String())
	defer span.End()

	currentPulse, err := b.PulseAccessor.Latest(ctx)
	if err != nil {
		return nil, err
	}

	parcel, err := b.CreateParcel(ctx, msg, ops.Safe().Token, currentPulse)
	if err != nil {
		return nil, err
	}

	rep, err := b.SendParcel(ctx, parcel, currentPulse, ops)
	return rep, err
}

// CreateParcel creates signed message from provided message.
func (b *Bus) CreateParcel(ctx context.Context, msg insolar.Message, token insolar.DelegationToken, currentPulse insolar.Pulse) (insolar.Parcel, error) {
	return b.ParcelFactory.Create(ctx, msg, b.NodeNetwork.GetOrigin().ID(), token, currentPulse)
}

// SendParcel sends provided message via network.
func (b *Bus) SendParcel(
	ctx context.Context,
	parcel insolar.Parcel,
	currentPulse insolar.Pulse,
	options *insolar.MessageSendOptions,
) (<-chan insolar.Reply, error) {
	f := newFuture()
	payload := insMsg.ParcelToBytes(parcel)
	msg := message.NewMessage(watermill.NewUUID(), payload)
	b.resultMutex.Lock()
	b.resultMap[msg.UUID] = f
	b.resultMutex.Unlock()
	msg.Metadata.Set("Pulse", fmt.Sprintf("%d", currentPulse.PulseNumber))
	msg.Metadata.Set("Type", parcel.Message().Type().String())
	fmt.Println("hi love, set type:", parcel.Message().Type().String(), msg.UUID)
	// b.Network.SendMessage()

	err := b.Pub.Publish("outbound", msg)
	if err != nil {
		// 	do staff
	}

	return f.Result(), nil
}
