/*
 *    Copyright 2019 Insolar
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

package docflow

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/sc"
	"testing"
)

type CallType string

const (
	QUERY  CallType = "query"
	INVOKE CallType = "invoke"
)

type Call struct {
	Name          string      `json:"name"`
	TransactionID string      `json:"transactionID"`
	CallType      CallType    `json:"callType"`
	In            interface{} `json:"in"`
	Out           interface{} `json:"out"`
}

type Emulation struct {
	Methods []Call `json:"methods"`
}

func ChatEmulation1(t *testing.T) *shim.MockStub {
	sc := new(sc.SmartContract)
	stub := shim.NewMockStub("guarantee_cc_test", sc)

	checkInit(t, stub, [][]byte{[]byte("init")})

	checkInvoke(t, stub, "11111", "initLedger", []string{})

	emulationJson := readJsonFile("emulation1")

	var emulation Emulation

	err := json.Unmarshal([]byte(emulationJson), &emulation)
	if err != nil {
		fmt.Println("Incorrect emulation json file. Error: ", err)
		t.FailNow()
	}
	fmt.Println(emulation)

	for _, m := range emulation.Methods {
		checkQuery(t, stub, m.TransactionID, m.Name, m.In, m.Out)
	}

	return stub
}

func Test_ChatEmulation(t *testing.T) {
	ChatEmulation1(t)
}
