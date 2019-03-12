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

package sc

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_sc"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
)

type SmartContract struct {
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) peer.Response {
	return com.SuccessMessageResponse("Starting chaincode!")
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) peer.Response {
	element := com.FPath.Path.PushBack("Invoke")
	defer com.FPath.Path.Remove(element)

	function, invokeArgs := APIstub.GetFunctionAndParameters()

	com.DebugLogMsg("Invoke of " + function)
	com.DebugLogMsg("Invoke args:  " + com.ConcatArrStr(invokeArgs))

	if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "hcInitLedger" {
		return s.hcInitLedger(APIstub)
	} else if function == "queryUser" {
		return s.queryUser(APIstub)
	}

	switch function {
	case "setDebugLogLevel":
		com.Logger.SetLevel(shim.LogDebug)
		return com.SuccessMessageResponse("Debug log level was set.")
	case "setInfoLogLevel":
		com.Logger.SetLevel(shim.LogInfo)
		return com.SuccessMessageResponse("Info log level was set.")
	case "setNoticeLogLevel":
		com.Logger.SetLevel(shim.LogNotice)
		return com.SuccessMessageResponse("Notice log level was set.")
	case "setWarningLogLevel":
		com.Logger.SetLevel(shim.LogWarning)
		return com.SuccessMessageResponse("Warning log level was set.")
	case "setErrorLogLevel":
		com.Logger.SetLevel(shim.LogError)
		return com.SuccessMessageResponse("Error log level was set.")
	case "setCriticalLogLevel":
		com.Logger.SetLevel(shim.LogCritical)
		return com.SuccessMessageResponse("Critical log level was set.")
	}

	return df_sc.CallBusinessFunc(APIstub, function, invokeArgs, "false")
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	com.Logger.SetLevel(shim.LogDebug)

	response := df_sc.CallBusinessFunc(APIstub, "hcCreateTemplateOfChat", []string{}, "false")
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return com.SuccessMessageResponse("Inited successfully!")
}

func (s *SmartContract) hcInitLedger(APIstub shim.ChaincodeStubInterface) peer.Response {
	com.Logger.SetLevel(shim.LogDebug)

	response := df_sc.CallBusinessFunc(APIstub, "hcCreateOrganizationsAndParticipants", []string{}, "false")
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	response = df_sc.CallBusinessFunc(APIstub, "hcCreateTemplateOfChat", []string{}, "false")
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	response = df_sc.CallBusinessFunc(APIstub, "hcCreateChat", []string{}, "false")
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return com.SuccessMessageResponse("Inited successfully!")
}
