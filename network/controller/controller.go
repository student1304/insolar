/*
 *    Copyright 2018 Insolar
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

package controller

import (
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/controller/auth"
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/insolar/insolar/network/transport/packet/types"
)

// Controller contains network logic.
type Controller struct {
	options common.Options
	network network.HostNetwork

	bootstrapController common.BootstrapController
	authController      *auth.AuthorizationController
	pulseController     *PulseController
}

// SendMessage send message to nodeID.
func (c *Controller) SendMessage(nodeID core.RecordRef, name string, msg core.SignedMessage) ([]byte, error) {
	return nil, nil
}

// RemoteProcedureRegister register remote procedure that will be executed when message is received.
func (c *Controller) RemoteProcedureRegister(name string, method core.RemoteProcedure) {

}

// SendCascadeMessage sends a message from MessageBus to a cascade of nodes.
func (c *Controller) SendCascadeMessage(data core.Cascade, method string, msg core.SignedMessage) error {
	return nil
}

// Bootstrap init bootstrap process: 1. Connect to discovery node; 2. Reconnect to new discovery node if redirected.
func (c *Controller) Bootstrap() error {
	return c.bootstrapController.Bootstrap()
}

// AnalyzeNetwork legacy method for old DHT network (should be removed in new network).
func (c *Controller) AnalyzeNetwork() error {
	log.Warn("this method was created for compatibility with old network, should be deleted")
	return nil
}

// Authorize start authorization process on discovery node.
func (c *Controller) Authorize() error {
	return c.authController.Authorize()
}

// ResendPulseToKnownHosts resend pulse when we receive pulse from pulsar daemon.
func (c *Controller) ResendPulseToKnownHosts(pulse core.Pulse) {
	c.pulseController.ResendPulse(pulse)
}

// GetNodeID get self node id (should be removed in far future).
func (c *Controller) GetNodeID() core.RecordRef {
	return core.RecordRef{}
}

// Inject inject components.
func (c *Controller) Inject(components core.Components) {
	c.network.RegisterRequestHandler(types.Ping, func(request network.Request) (network.Response, error) {
		return c.network.BuildResponse(request, nil), nil
	})
	c.bootstrapController.Start()
	c.authController.Start(components)
	c.pulseController.Start()
}

// NewNetworkController create new network controller.
func NewNetworkController(
	pulseCallback network.OnPulse,
	configuration configuration.Configuration,
	transport hostnetwork.InternalTransport,
	routingTable network.RoutingTable,
	network network.HostNetwork) network.Controller {

	c := Controller{}
	c.network = network
	c.bootstrapController = NewBootstrapController(&c.options, transport)
	c.authController = auth.NewAuthorizationController(&c.options, c.bootstrapController, transport)
	c.pulseController = NewPulseController(pulseCallback, network, routingTable)

	return &c
}
