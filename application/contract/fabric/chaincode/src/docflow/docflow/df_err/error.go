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

package df_err

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
)

const (
	NO_PERMISION_FOR_WORK_WITH_CHAT_ERROR    = 1010
	NO_PARTICIPANT_FOR_EXECUTOR_ERROR        = 1020
	EMPTY_PROC_TEMPLATE_ELEMENTS_ERROR       = 1030
	HASH_IS_ALREADY_SET_ERROR                = 1040
	ATTACHMENT_WITH_THIS_NAME_ALREADY_EXISTS = 1040
)

func NoParticipantForExecutor(executorMSPId string) peer.Response {
	com.DebugLogMsg("Executors MSPId: " + executorMSPId)
	message := "No permission for work in flow."
	return com.ErrorMessageResponse(com.CCError{"NoParticipantForExecutor"}, NO_PARTICIPANT_FOR_EXECUTOR_ERROR, message)
}

func NoPermissionForWorkWithProcessError(executorParticipantReference string) peer.Response {
	com.DebugLogMsg("Your participant reference: " + executorParticipantReference)
	message := "No permission for work with process."
	return com.ErrorMessageResponse(com.CCError{"NoPermissionForWorkWithProcessError"}, NO_PERMISION_FOR_WORK_WITH_CHAT_ERROR, message)
}

func EmptyProcTemplateElementsError() peer.Response {
	message := "Process template elements count is less than 3."
	return com.ErrorMessageResponse(com.CCError{"EmptyProcTemplateElementsError"}, EMPTY_PROC_TEMPLATE_ELEMENTS_ERROR, message)
}

func HashIsAlreadySetError() peer.Response {
	message := "Hash is already set."
	return com.ErrorMessageResponse(com.CCError{"HashIsAlreadySetError"}, HASH_IS_ALREADY_SET_ERROR, message)
}

func AttachmentWithThisNameAlreadyExistsError(name string) peer.Response {
	message := "Attachment with name <<" + name + ">> already exists."
	return com.ErrorMessageResponse(com.CCError{"AttachmentWithThisNameAlreadyExistsError"}, ATTACHMENT_WITH_THIS_NAME_ALREADY_EXISTS, message)
}
