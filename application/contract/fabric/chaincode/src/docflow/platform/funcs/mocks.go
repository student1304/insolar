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

package funcs

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
	"strconv"
)

const (
	MSPId     = "Org1MSP"
	Timestamp = 1549620261
)

//func GetMSPID(stub shim.ChaincodeStubInterface) (string, peer.Response) {
//	element := com.FPath.Path.PushBack("GetMSPID")
//	defer com.FPath.Path.Remove(element)
//
//	mspId, err := cid.GetMSPID(stub)
//	if err != nil {
//		return "", com.GetMSPIDError(err)
//	}
//	return mspId, com.SuccessMessageResponse("MSP Id was gotten.")
//}
//
//func GetTime(stub shim.ChaincodeStubInterface) (int64, peer.Response) {
//	element := com.FPath.Path.PushBack("GetTime")
//	defer com.FPath.Path.Remove(element)
//
//	time, err := stub.GetTxTimestamp()
//	if err != nil {
//		return 0, com.GetTxTimestampError(err)
//	}
//	return time.Seconds, com.SuccessMessageResponse("Transaction timestamp was gotten.")
//}

func GetMSPID(stub shim.ChaincodeStubInterface) (string, peer.Response) {
	element := com.FPath.Path.PushBack("GetMSPID")
	defer com.FPath.Path.Remove(element)

	return MSPId, com.SuccessMessageResponse("Mock MSP_Id <<" + MSPId + ">> was gotten.")
}

func GetTime(stub shim.ChaincodeStubInterface) (int64, peer.Response) {
	element := com.FPath.Path.PushBack("GetTime")
	defer com.FPath.Path.Remove(element)

	return Timestamp, com.SuccessMessageResponse("Mock timestamp <<" + strconv.Itoa(Timestamp) + ">> was gotten.")
}
