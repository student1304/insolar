//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.

package gateway

import (
	"context"

	"github.com/insolar/insolar/log" // TODO remove before merge

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network"
)

// Base is abstract class for gateways

type Base struct {
	Me                 network.Gateway
	Network            network.Gatewayer
	GIL                insolar.GlobalInsolarLock
	SwitcherWorkAround insolar.SwitcherWorkAround // nodekeeper
}

// NewGateway creates new gateway on top of existing
func (g *Base) NewGateway(state insolar.NetworkState) network.Gateway {
	log.Warnf("NewGateway %s", state.String())
	switch state {
	case insolar.NoNetworkState:
		g.Me = &NoNetwork{g}
	case insolar.VoidNetworkState:
		g.Me = NewVoid(g)
	case insolar.JetlessNetworkState:
		g.Me = NewJetless(g)
	case insolar.AuthorizationNetworkState:
		g.Me = NewAuthorisation(g)
	case insolar.CompleteNetworkState:
		g.Me = NewComplete(g)
	default:
		panic("Try to switch network to unknown state. Memory of process is inconsistent.")
	}
	return g.Me
}

func (g *Base) OnPulse(ctx context.Context, pu insolar.Pulse) error {
	if g.SwitcherWorkAround.IsBootstrapped() {
		g.Network.SetGateway(g.Network.Gateway().NewGateway(insolar.CompleteNetworkState))
		g.Network.Gateway().Run(ctx)
	}
	return nil
}

// Auther casts us to Auther or obtain it in another way
func (g *Base) Auther() network.Auther {
	if ret, ok := g.Me.(network.Auther); ok {
		return ret
	}
	panic("Our network gateway suddenly is not an Auther")
}

// GetCert method returns node certificate by requesting sign from discovery nodes
func (nc *Base) GetCert(ctx context.Context, ref *insolar.Reference) (insolar.Certificate, error) {
	panic("GetCert() is not useable in this state")
}

// ValidateCert validates node certificate
func (nc *Base) ValidateCert(ctx context.Context, certificate insolar.AuthorizationCertificate) (bool, error) {
	panic("ValidateCert()  is not useable in this state")
}