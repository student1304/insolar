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

package com

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
)

var Logger = shim.NewLogger("myChaincode")

func InfoLogMsg(message string) {
	Logger.Info("Message: ", FPath.getPath()+message)
}

func DebugLogMsg(message string) {
	Logger.Debug("Message: ", FPath.getPath()+message)
}

func ErrorLogMsg(error error, status int32, message string) {
	Logger.Error("Code: ", strconv.Itoa(int(status)),
		"; Message: ", message,
		"; Error: ", error.Error())
}
