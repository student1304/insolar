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
	"github.com/hyperledger/fabric/protos/peer"
	"strconv"
)

const (
	INVALID_FUNCTION_NAME_ERROR              = 510
	INCORRECT_INVOKE_NUMBER_OF_ARGS_ERROR    = 511
	INCORRECT_INVOKE_ARGS_ERROR              = 512
	INCORRECT_NUMBER_OF_ARGS_ERROR           = 512
	MARSHAL_ERROR                            = 520
	RESPONSE_MARSHAL_ERROR                   = 521
	UNMARSHAL_ERROR                          = 530
	GET_STATE_ERROR                          = 540
	GET_STATE_BY_PARTIAL_COMPOSITE_KEY_ERROR = 541
	PUT_STATE_ERROR                          = 550
	PUT_STATE_ON_EXIST_KEY_ERROR             = 551
	EDIT_STATE_ON_NOT_EXIST_KEY_ERROR        = 552
	CREATE_COMPOSITE_KEY_ERROR               = 560
	SPLIT_COMPOSITE_KEY_ERROR                = 561
	ITERATOR_NEXT_ERROR                      = 570
	ENCRYPT_ERROR                            = 580
	DECRYPT_ERROR                            = 590
	DECODE_ERROR                             = 591
	GET_ATTRIBUTE_ERROR                      = 600
	NOT_POSSESS_ATTRIBUTE_ERROR              = 601
	GET_MSP_ID_ERROR                         = 610
	INCORRECT_NUMBER_OF_FIELDS_ERROR         = 620
	NO_SUCH_CASE_OF_ENTITY_ERROR             = 630
	INCORRECT_ARRAY_PARSE_ARGS_ERROR         = 640
	ENTITY_VALIDATION_ERROR                  = 650
	EVENT_ERROR                              = 660
	GET_TX_TIMESTAMP_ERROR                   = 670
)

type CCError struct {
	Name string `protobuf:"bytes,1,opt,name=Name" json:"Name" json:"Name"`
}

func (error CCError) Error() string {
	return error.Name
}

func InvalidFunctionNameForRoleError(function string, role string) peer.Response {
	message := "Invalid Smart Contract function name: '" + function + "' for role: '" + role + "'"
	return ErrorMessageResponse(CCError{"InvalidFunctionNameForRoleError"}, INVALID_FUNCTION_NAME_ERROR, message)
}

func IncorrectInvokeNumberOfArgsError(args []string, narg int) peer.Response {
	message := "Incorrect invoke arguments: " + ConcatArrStr(args) + ". Actual number of args: " + strconv.Itoa(len(args)) + " Expecting number of args: " + strconv.Itoa(narg)
	return ErrorMessageResponse(CCError{"IncorrectInvokeNumberOfArgsError"}, INCORRECT_INVOKE_NUMBER_OF_ARGS_ERROR, message)
}

func IncorrectInvokeArgsError(args []string) peer.Response {
	message := "Incorrect invoke arguments: " + ConcatArrStr(args)
	return ErrorMessageResponse(CCError{"IncorrectInvokeArgsError"}, INCORRECT_INVOKE_ARGS_ERROR, message)
}

func IncorrectNumberOfArgsError(args map[string]string, narg int) peer.Response {
	message := "Incorrect number of arguments. Expecting " + strconv.Itoa(narg)
	return ErrorMessageResponse(CCError{"IncorrectNumberOfArgsError"}, INCORRECT_NUMBER_OF_ARGS_ERROR, message)
}

func MarshalError(err error) peer.Response {
	message := "Error on Marshal data"
	return ErrorMessageResponse(err, MARSHAL_ERROR, message)
}

func ResponseMarshalError(err error) peer.Response {
	message := "Error on Marshal response"
	return ErrorMessageResponse(err, RESPONSE_MARSHAL_ERROR, message)
}

func UnmarshalError(err error, data string) peer.Response {
	message := "Error on Unmarshal data: " + data
	return ErrorMessageResponse(err, UNMARSHAL_ERROR, message)
}

func GetStateError(err error, data string) peer.Response {
	message := "Error on GetState: " + data
	return ErrorMessageResponse(err, GET_STATE_ERROR, message)
}

func GetStateByPartialCompositeKeyError(err error, objectType string, keys []string) peer.Response {
	message := "Error on GetStateByPartialCompositeKey for objectType = " + objectType + "; keys = " + ConcatArrStr(keys)
	return ErrorMessageResponse(err, GET_STATE_BY_PARTIAL_COMPOSITE_KEY_ERROR, message)
}

func PutStateError(err error, key string, data []byte) peer.Response {
	str := "Key: " + key + "; data: " + string(data)
	message := "Error on PutState: " + str
	return ErrorMessageResponse(err, PUT_STATE_ERROR, message)
}

func PutStateOnExistKeyError(key string, data []byte) peer.Response {
	str := "Key: " + key + "; data: " + string(data)
	message := "Error on PutState, the key already exist! " + str
	return ErrorMessageResponse(CCError{"PutStateOnExistKeyError"}, PUT_STATE_ON_EXIST_KEY_ERROR, message)
}

func EditStateOnNotExistKeyError(key string, data []byte) peer.Response {
	str := "Key: " + key + "; data: " + string(data)
	message := "Error on PutState, the key does not exist! " + str
	return ErrorMessageResponse(CCError{"EditStateOnNotExistKeyError"}, EDIT_STATE_ON_NOT_EXIST_KEY_ERROR, message)
}

func CreateCompositeKeyError(err error, objectType string, attributes []string) peer.Response {
	data := "objectType: " + objectType + "; attributes: " + ConcatArrStr(attributes)
	message := "Error on CreateCompositeKey: " + data
	return ErrorMessageResponse(err, CREATE_COMPOSITE_KEY_ERROR, message)
}

func SplitCompositeKeyError(err error, compositeKey string) peer.Response {
	data := "compositeKey: " + compositeKey
	message := "Error on SplitCompositeKey: " + data
	return ErrorMessageResponse(err, SPLIT_COMPOSITE_KEY_ERROR, message)
}

func IteratorNextError(err error) peer.Response {
	message := "Error on get next item from iterator."
	return ErrorMessageResponse(err, ITERATOR_NEXT_ERROR, message)
}

func EncryptError(err error, data string) peer.Response {
	message := "Error on encrypt: " + data
	return ErrorMessageResponse(err, ENCRYPT_ERROR, message)
}

func DecryptError(err error, data string) peer.Response {
	message := "Error on decrypt: " + data
	return ErrorMessageResponse(err, DECRYPT_ERROR, message)
}

func DecodeError(err error, data string) peer.Response {
	message := "Error on decode: " + data
	return ErrorMessageResponse(err, DECODE_ERROR, message)
}

func GetAttributeError(err error, data string) peer.Response {
	message := "Error on get attribute: " + data
	return ErrorMessageResponse(err, GET_ATTRIBUTE_ERROR, message)
}

func NotPossessAttributeError(data string) peer.Response {
	message := "Error on getting attribute. Not possess attribute: " + data
	return ErrorMessageResponse(CCError{"NotPossessAttributeError"}, NOT_POSSESS_ATTRIBUTE_ERROR, message)
}

func GetMSPIDError(err error) peer.Response {
	message := "Error on getting MSPId!"
	return ErrorMessageResponse(err, GET_MSP_ID_ERROR, message)
}

func IncorectNumberOfFieldsError(nFieldPaths, nFieldValues int) peer.Response {
	message := "Incorrect number of fields values and fields paths. They are not equal. FieldPath paths: " +
		strconv.Itoa(nFieldPaths) + "; FieldPath values: " + strconv.Itoa(nFieldValues)
	return ErrorMessageResponse(CCError{"IncorrectNumberOfArgsError"}, INCORRECT_NUMBER_OF_FIELDS_ERROR, message)
}

func NoSuchCaseOfEntityError(entityName string) peer.Response {
	message := "No such case of entity. Entity name: " + entityName
	return ErrorMessageResponse(CCError{"NoSuchCaseOfEntityError"}, NO_SUCH_CASE_OF_ENTITY_ERROR, message)
}

func IncorrectArrayParseArgsError() peer.Response {
	message := "Incorrect array parse arguments."
	return ErrorMessageResponse(CCError{"NoSuchCaseOfEntityError"}, INCORRECT_ARRAY_PARSE_ARGS_ERROR, message)
}

func EntityValidationError() peer.Response {
	message := "Cann't edit entity because no permission for editing some field. "
	return ErrorMessageResponse(CCError{"EntityValidationError"}, ENTITY_VALIDATION_ERROR, message)
}

func EventError(err error, eventName string) peer.Response {
	return ErrorMessageResponse(err, EVENT_ERROR, "Event "+eventName+" error.")
}

func GetTxTimestampError(err error) peer.Response {
	return ErrorMessageResponse(err, GET_TX_TIMESTAMP_ERROR, "Transaction timestamp getting error")
}

func ConcatArrStr(arr []string) string {
	result := "["
	for index, element := range arr {
		result += "Index " + strconv.Itoa(index) + ", Value " + element + ";"
	}
	result += "]"

	return result
}
