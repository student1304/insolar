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
	"encoding/base64"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
)

func IdToKey(id string) (string, peer.Response) {
	element := com.FPath.Path.PushBack("IdToKey")
	defer com.FPath.Path.Remove(element)

	result, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		com.DecodeError(err, "Data: "+id)
	}

	return string(result), com.SuccessMessageResponse("String was decrypted successfully.")
}

func KeyToId(key string) (id string) {
	return base64.StdEncoding.EncodeToString([]byte(key))
}
