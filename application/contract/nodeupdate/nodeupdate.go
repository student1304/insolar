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

package nodeupdate

import (
	"time"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// NodeUpdate is smart contract representing cloud update control logic
type NodeUpdate struct {
	foundation.BaseContract
	UpdateTimeLimit   int64
	AllowedUpdateRate float32
	NodeCount         uint32
	mandates          map[core.RecordRef]int64
}

// RequestMandate allows or denies update request
func (nu *NodeUpdate) RequestMandate() bool {
	ctx := nu.GetContext()
	caller := *ctx.Caller
	updateRate := float32(len(nu.mandates)) / float32(nu.NodeCount)
	if updateRate < nu.AllowedUpdateRate {
		requestTime := time.Now().Unix()
		nu.mandates[caller] = requestTime
		return true
	}
	return false
}

// ReleaseMandate signals to the cloud that the node has been updated
func (nu *NodeUpdate) ReleaseMandate() {
	ctx := nu.GetContext()
	caller := *ctx.Caller
	delete(nu.mandates, caller)
}

// CheckMandate allows the node to know if the update is still possible
func (nu *NodeUpdate) CheckMandate() bool {
	ctx := nu.GetContext()
	caller := *ctx.Caller
	startTime := nu.mandates[caller]
	if startTime > 0 {
		now := time.Now().Unix()
		if now > startTime+nu.UpdateTimeLimit {
			delete(nu.mandates, caller)
			return false
		}
		return true
	}
	return false
}

// Cleanup releases overdue update permissions
func (nu *NodeUpdate) Cleanup() {
	now := time.Now().Unix()
	for nodeRef, startTime := range nu.mandates {
		if now > startTime+nu.UpdateTimeLimit {
			delete(nu.mandates, nodeRef)
		}
	}
}

// NewNodeUpdate creates new NodeUpdate
func NewNodeUpdate(timeLimit int64, updateRate float32, nodeCount uint32) (*NodeUpdate, error) {
	return &NodeUpdate{
		UpdateTimeLimit:   timeLimit,
		AllowedUpdateRate: updateRate,
		NodeCount:         nodeCount,
	}, nil
}
