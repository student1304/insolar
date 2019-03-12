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

package user

import (
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
)

// Define  structure.
type User struct {
	FIO          string `json:"fio"`
	Position     string `json:"position"`
	Role         string `json:"role"`
	CreationDate string `json:"creation_date"`
	CreatorLogin string `json:"creator_login"`
}

const (
	FIO          = "fio"
	POSITION     = "position"
	ROLE         = "role"
	CREATIONDATE = "creation_date"
	CREATORLOGIN = "creator_login"

	ROLE_READER  = "reader"
	ROLE_CREATER = "creator"
	ROLE_ADMIN   = "admin"
	ROLE_BANK    = "bank"
)

func GetUser(APIstub shim.ChaincodeStubInterface) (User, peer.Response) {
	element := com.FPath.Path.PushBack("GetUser")
	defer com.FPath.Path.Remove(element)

	user := User{}

	val, ok, err := cid.GetAttributeValue(APIstub, FIO)
	if err != nil {
		return user, com.GetAttributeError(err, FIO)
	}
	if !ok {
		return user, com.NotPossessAttributeError(FIO)
	}
	user.FIO = val

	val, ok, err = cid.GetAttributeValue(APIstub, POSITION)
	if err != nil {
		return user, com.GetAttributeError(err, POSITION)
	}
	if !ok {
		return user, com.NotPossessAttributeError(POSITION)
	}
	user.Position = val

	val, ok, err = cid.GetAttributeValue(APIstub, ROLE)
	if err != nil {
		return user, com.GetAttributeError(err, ROLE)
	}
	if !ok {
		return user, com.NotPossessAttributeError(ROLE)
	}
	user.Role = val

	val, ok, err = cid.GetAttributeValue(APIstub, CREATIONDATE)
	if err != nil {
		return user, com.GetAttributeError(err, CREATIONDATE)
	}
	if !ok {
		return user, com.NotPossessAttributeError(CREATIONDATE)
	}
	user.CreationDate = val

	val, ok, err = cid.GetAttributeValue(APIstub, CREATORLOGIN)
	if err != nil {
		return user, com.GetAttributeError(err, CREATORLOGIN)
	}
	if !ok {
		return user, com.NotPossessAttributeError(CREATORLOGIN)
	}
	user.CreatorLogin = val

	return user, com.SuccessMessageResponse("User was gotten from attributes.")
}
