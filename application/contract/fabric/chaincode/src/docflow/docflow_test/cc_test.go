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
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"unicode"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}

	checkInvoke(t, stub, "1", "setDebugLogLevel", []string{})
	//checkInvoke(t, stub, "1", "setInfoLogLevel", []string{})
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
}

func checkQuery(t *testing.T, stub *shim.MockStub, txId, functionName string, inputJson, outputJson interface{}) {

	inputAsBytes, err := json.Marshal(inputJson)
	if err != nil {
		fmt.Println("Error on Marshal input json.")
		t.FailNow()
	}
	outputAsBytes, err := json.Marshal(outputJson)
	if err != nil {
		fmt.Println("Error on Marshal output json.")
		t.FailNow()
	}

	argsAsBytes := [][]byte{[]byte(functionName), inputAsBytes}

	res := stub.MockInvoke(txId, argsAsBytes)

	removeSpacesValue := removeSpaces(string(outputAsBytes))
	removeSpacesPayload := removeSpaces(string(res.Payload))
	if removeSpaces(removeSpacesPayload) != removeSpacesValue {
		fmt.Println("\n"+functionName+" value \n", removeSpacesPayload, "\n was not\n", removeSpacesValue, "\n as expected ")
		t.FailNow()
	} else {
		fmt.Println("\nQueried ", functionName, " successfully!")
	}
}

func checkInvoke(t *testing.T, stub *shim.MockStub, txId, functionName string, args []string) {
	argsAsBytes := [][]byte{[]byte(functionName)}
	for _, arg := range args {
		argsAsBytes = append(argsAsBytes, []byte(arg))
	}

	res := stub.MockInvoke(txId, argsAsBytes)
	if res.Status != shim.OK {
		fmt.Println("\nInvoke\n", args, "\n failed\n", string(res.Message))
		t.FailNow()
	} else {
		fmt.Println("\nInvoked ", functionName, " successfully!")
	}
}

func readJsonFile(fileName string) string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Print(err)
	}

	fileContent, err := ioutil.ReadFile(dir + "/testdata/" + fileName + ".json")
	if err != nil {
		fmt.Print(err)
	}

	return string(fileContent)
}

func removeSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
