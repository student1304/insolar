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
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/login/user"
)

func (s *SmartContract) queryUser(APIstub shim.ChaincodeStubInterface) peer.Response {
	element := com.FPath.Path.PushBack("s.queryUser")
	defer com.FPath.Path.Remove(element)

	user, response := user.GetUser(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return com.SuccessPayloadResponse(&user)
}
