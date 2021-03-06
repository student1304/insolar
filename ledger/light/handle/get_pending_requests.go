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

package handle

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/flow/bus"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type GetPendingRequests struct {
	dep      *proc.Dependencies
	msg      *message.GetPendingRequests
	replyTo  chan<- bus.Reply
	reqPulse insolar.PulseNumber
}

func NewGetPendingRequests(dep *proc.Dependencies, rep chan<- bus.Reply, parcel insolar.Parcel) *GetPendingRequests {
	return &GetPendingRequests{
		dep:      dep,
		msg:      parcel.Message().(*message.GetPendingRequests),
		replyTo:  rep,
		reqPulse: parcel.Pulse(),
	}
}

func (s *GetPendingRequests) Present(ctx context.Context, f flow.Flow) error {
	jet := proc.NewFetchJet(*s.msg.DefaultTarget().Record(), flow.Pulse(ctx), s.replyTo)
	s.dep.FetchJet(jet)
	if err := f.Procedure(ctx, jet, false); err != nil {
		return err
	}

	getPendingRequests := proc.NewGetPendingRequests(jet.Result.Jet, s.replyTo, s.msg, s.reqPulse)
	s.dep.GetPendingRequests(getPendingRequests)
	return f.Procedure(ctx, getPendingRequests, false)
}
