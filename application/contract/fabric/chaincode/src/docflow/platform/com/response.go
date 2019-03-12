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
	"container/list"
	"encoding/json"
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

const (
	// OK_MESSAGE constant - status code less than 400, endorser will endorse it.
	// OK_MESSAGE means init or invoke successfully.
	OK         = 200
	OK_MESSAGE = 0
	OK_PAYLOAD = 0

	// ERRORTHRESHOLD constant - status code greater than or equal to 400 will be considered an Name and rejected by endorser.
	ERRORTHRESHOLD = 400

	// ERROR constant - default Name value
	ERROR = 500
)

var FPath_obj = FuncPath{Path: list.New()}
var FPath *FuncPath = &FPath_obj

type FuncPath struct {
	Path *list.List `protobuf:"bytes,1,opt,name=path" json:"path" xml:"path"`
}

type Response struct {
	// A status code that should follow the HTTP status codes.
	Status int32 `protobuf:"varint,1,opt,name=status" json:"status,omitempty" xml:"status"`
	// A Name message that may be kept.
	Error string `protobuf:"varint,1,opt,name=Name" json:"Name,omitempty" xml:"Name,omitempty"`
	// A message associated with the response code.
	Message string `protobuf:"bytes,2,opt,name=message" json:"message,omitempty" xml:"message"`
	// A payload that can be used to include metadata with this response.
	Payload []byte `protobuf:"bytes,3,opt,name=payload,proto3" json:"payload,omitempty" xml:"payload"`
}

func (FPath FuncPath) getPath() string {
	if FPath.Path == nil {
		return ""
	}

	result := " '"

	for e := FPath.Path.Front(); e != nil; e = e.Next() {
		result += "/" + (e.Value).(string)
	}

	result += "': "

	return result
}

func ErrorMessageResponse(error error, status int32, message string) peer.Response {
	ErrorLogMsg(error, status, message)

	return peer.Response{
		Status:  ERROR,
		Message: " Error code: " + strconv.Itoa(int(status)) + ": " + error.Error() + ". Message: " + message,
	}
}

func SuccessMessageResponse(message string) peer.Response {
	DebugLogMsg(message)

	return peer.Response{
		Status:  OK,
		Message: message,
	}
}

func SuccessPayloadResponse(data interface{}) peer.Response {

	dataJSON, err := json.Marshal(data)
	if err != nil {
		ResponseMarshalError(err)
	}

	DebugLogMsg(string(dataJSON))

	return peer.Response{
		Status:  OK,
		Payload: dataJSON,
	}
}
