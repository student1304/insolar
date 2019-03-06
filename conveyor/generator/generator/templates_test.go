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

package generator

const (
	genStateMachine = `
package sample

import (
	"github.com/insolar/insolar/conveyor/generator/common"
	"errors"
)

type SMFIDTestStateMachine struct { // STFID = State Machine Flow IDs
}

func (*SMFIDTestStateMachine) StateFirst() common.ElState {
    return 1
}

func (*SMFIDTestStateMachine) StateSecond() common.ElState {
    return 2
}
`
	genRawHandlers = `
type SMRHTestStateMachine struct { // SMRH = State Machine Raw Handlers
	cleanHandlers TestStateMachine
}

func NewSMRHTestStateMachine() SMRHTestStateMachine {
	return SMRHTestStateMachine{
		// cleanHandlers: &TestStateMachineImplementation{},
	}
}

func (s *SMRHTestStateMachine) Init(element slot.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    payload, state, err := s.cleanHandlers.Init(aInput)
    if err != nil {
        return payload, state, err
    }
    return s.cleanHandlers.TransitFirstSecond(aInput, payload)
}

func (s *SMRHTestStateMachine) TransitFirstSecond(element slot.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.TransitFirstSecond(aInput, aPayload)
}

func (s *SMRHTestStateMachine) MigrateFirst(element slot.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.MigrateFirst(aInput, aPayload)
}

func (s *SMRHTestStateMachine) ErrorFirst(element slot.SlotElementHelper, err error) (interface{}, common.ElState) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong input event type")
        return nil, 0
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong payload type")
        return nil, 0
    }
    return s.cleanHandlers.ErrorFirst(aInput, aPayload, err)
}

func (s *SMRHTestStateMachine) TransitSecondThird(element slot.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.TransitSecondThird(aInput, aPayload)
}

func (s *SMRHTestStateMachine) MigrateSecond(element slot.SlotElementHelper) (interface{}, common.ElState, error) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        return nil, 0, errors.New("wrong input event type")
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        return nil, 0, errors.New("wrong payload type")
    }
    return s.cleanHandlers.MigrateSecond(aInput, aPayload)
}

func (s *SMRHTestStateMachine) ErrorSecond(element slot.SlotElementHelper, err error) (interface{}, common.ElState) {
    aInput, ok := element.GetInputEvent().(Event)
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong input event type")
        return nil, 0
    }
    aPayload, ok := element.GetPayload().(*Payload)
    if !ok {
        // TODO fix me
        // return nil, 0, errors.New("wrong payload type")
        return nil, 0
    }
    return s.cleanHandlers.ErrorSecond(aInput, aPayload, err)
}
`
)

// func TestGenerator_GenerateStateMachine(t *testing.T) {
// 	g := testGenerator(t)
// 	g.findEachStateMachine()
// 	out := new(bytes.Buffer)
// 	g.GenerateStateMachine(out, 0)
// 	assert.Equal(t, genStateMachine, out.String())
// 	fmt.Println(out.String())
// }
//
// func TestGenerator_GenerateRawHandlers(t *testing.T) {
// 	g := testGenerator(t)
// 	g.findEachStateMachine()
// 	out := new(bytes.Buffer)
// 	g.GenerateRawHandlers(out, 0)
// 	assert.Equal(t, genRawHandlers, out.String())
// }