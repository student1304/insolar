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
	"encoding/json"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
)

func IsFieldsEqual(field1, field2 interface{}) bool {
	fieldBytes1, err := json.Marshal(field1)
	if err != nil {
		com.DebugLogMsg("Cann't marshal field 1.")
		return false
	}
	fieldBytes2, err := json.Marshal(field1)
	if err != nil {
		com.DebugLogMsg("Cann't marshal field 2.")
		return false
	}

	if string(fieldBytes1) != string(fieldBytes2) {
		com.DebugLogMsg("Field1 is not equal field2. Field1: " + string(fieldBytes1) + ". Field2: " + string(fieldBytes2))
		return false
	}

	return true
}
