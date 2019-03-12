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

package df_sc

import (
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/attachment"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/part"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/process"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/proctemplate"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/version"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/env/participant"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_err"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_templates"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/com"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/data"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/funcs"
)

func CallBusinessFunc(APIstub shim.ChaincodeStubInterface, function string, args []string, simulate string) peer.Response {
	switch function {

	case "createParticipant":
		return createParticipant(APIstub, args, simulate)
	case "editParticipant":
		return editParticipant(APIstub, args, simulate)
	case "getMyParticipants":
		return getMyParticipants(APIstub)
	case "getAllParticipants":
		return getAllParticipants(APIstub)
	case "hcCreateTemplateOfChat":
		return hcCreateTemplateOfChat(APIstub, simulate)

	case "createProcess":
		return createProcess(APIstub, args, simulate)
	case "getAttachmentsByProcessReference":
		return getAttachmentsByProcessReference(APIstub, args)
	case "getProcessElementHistory":
		return getProcessElementHistory(APIstub, args)
	case "getAllProcesses":
		return getAllProcesses(APIstub)
	case "createAttachment":
		return createAttachment(APIstub, args, simulate)
	case "createPart":
		return createPart(APIstub, args, simulate)
	case "sendResponse":
		return sendResponse(APIstub, args, simulate)
	case "createAttachmentVersion":
		return createAttachmentVersion(APIstub, args, simulate)
	case "getPartsByVersionReference":
		return getPartsByVersionReference(APIstub, args)
	case "getPartByReference":
		return getPartByReference(APIstub, args)

	default:
		return com.InvalidFunctionNameForRoleError(function, "")
	}
}

func createParticipant(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createParticipant")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	participantJSON := args[0]

	var participantObject participant.Participant

	err := json.Unmarshal([]byte(participantJSON), &participantObject)
	if err != nil {
		return com.UnmarshalError(err, participantJSON)
	}

	return data.Put(&participantObject, APIstub, simulate)
}

func editParticipant(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("editParticipant")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	participantJSON := args[0]

	var participantObject participant.Participant

	err := json.Unmarshal([]byte(participantJSON), &participantObject)
	if err != nil {
		return com.UnmarshalError(err, participantJSON)
	}

	return data.EditAll(&participantObject, APIstub, simulate)
}

func getMyParticipants(APIstub shim.ChaincodeStubInterface) peer.Response {
	element := com.FPath.Path.PushBack("getMyParticipants")
	defer com.FPath.Path.Remove(element)

	mspID, response := funcs.GetMSPID(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	var participantObjectsEntities []data.Entity
	data.QueryByIndex(&participant.Participant{}, &participantObjectsEntities, APIstub, []string{"MspId"}, mspID)

	return com.SuccessPayloadResponse(data.EntitiesToOut(&participant.Participant{}, participantObjectsEntities))
}

func getAllParticipants(APIstub shim.ChaincodeStubInterface) peer.Response {
	element := com.FPath.Path.PushBack("getAllParticipants")
	defer com.FPath.Path.Remove(element)

	var participantObjectsEntities []data.Entity
	data.QueryAll(&participant.Participant{}, &participantObjectsEntities, APIstub)

	return com.SuccessPayloadResponse(data.EntitiesToOut(&participant.Participant{}, participantObjectsEntities))
}

func hcCreateTemplateOfChat(APIstub shim.ChaincodeStubInterface, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("hcCreateTemplateOfChat")
	defer com.FPath.Path.Remove(element)

	return data.Put(&df_templates.ChatTemplate, APIstub, simulate)
}

func createProcTemplate(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createProcTemplate")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	procTemplateJsonStr := args[0]

	var procTemplateObject proctemplate.ProcTemplate

	err := json.Unmarshal([]byte(procTemplateJsonStr), &procTemplateObject)
	if err != nil {
		return com.UnmarshalError(err, procTemplateJsonStr)
	}

	if len(procTemplateObject.NcProcTemplate.Elements) < 3 {
		return df_err.EmptyProcTemplateElementsError()
	}

	return data.Put(&procTemplateObject, APIstub, simulate)
}

func getExecutorFirstParticipantRef(APIstub shim.ChaincodeStubInterface) (string, peer.Response) {
	element := com.FPath.Path.PushBack("getExecutorFirstParticipantRef")
	defer com.FPath.Path.Remove(element)

	mspID, response := funcs.GetMSPID(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return "", response
	}

	var participantObjectsEntities []data.Entity
	data.QueryByIndex(&participant.Participant{}, &participantObjectsEntities, APIstub, []string{"MspId"}, mspID)

	if len(participantObjectsEntities) == 0 {
		return "", df_err.NoParticipantForExecutor(mspID)
	}

	return participantObjectsEntities[0].GetId(), com.SuccessMessageResponse("First participant reference was gotten.")

}

func createProcess(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createProcess")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		Name                  string   `json:"name"`
		ProcTemplateReference string   `json:"processTemplateReference"`
		Participants          []string `json:"participants"` // not_docflow
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	procTemplateObject := proctemplate.ProcTemplate{Id: request.ProcTemplateReference}
	response := data.QueryById(&procTemplateObject, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	ref, response := getExecutorFirstParticipantRef(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	time, response := funcs.GetTime(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	processObject := process.Process{
		Name:           request.Name,
		Creator:        ref,
		CreationTime:   time,
		NcProcTemplate: procTemplateObject.NcProcTemplate,
	}

	// Start Set elements approve maps for process

	// docflow begin /////////////////
	//if len(processObject.NcProcTemplate.Elements)-2 != len(request.Participants) {
	//	return com.IncorrectInvokeArgsError(args)
	//}
	// docflow end /////////////////

	if len(processObject.NcProcTemplate.Elements) < 3 {
		return df_err.EmptyProcTemplateElementsError()
	}
	// docflow begin /////////////////
	//for i := 2; i < len(processObject.NcProcTemplate.Elements); i++ {
	//	for j := 0; j < len(request.Participants[i-2]); j++ {
	//		processObject.NcProcTemplate.Elements[i].ParticipantsApproves[request.Participants[i-2][j]] = false
	//	}
	//}
	// docflow end /////////////////

	// not_docflow
	for j := 0; j < len(request.Participants); j++ {
		processObject.NcProcTemplate.Elements[2].ParticipantsApproves[request.Participants[j]] = proctemplate.NONE
	}
	// End Set

	processObject.NcProcTemplate.Elements[0].Active = false
	processObject.NcProcTemplate.Elements[2].Active = true

	response = data.Put(&processObject, APIstub, simulate)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	// Begin event
	response = addParticipantsArrToProcessObject(APIstub, &processObject)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	eventJson, err := json.Marshal(processObject.ToOut())
	if err != nil {
		return com.MarshalError(err)
	}
	eventError := APIstub.SetEvent("new_process", eventJson)
	if err != nil {
		return com.EventError(eventError, "new_process")
	}
	// End event

	type OutputResponse struct {
		CreationDate     int64  `json:"creationDate"`
		ProcessReference string `json:"processReference"`
	}

	outputResponse := OutputResponse{
		ProcessReference: processObject.Id,
		CreationDate:     processObject.CreationTime,
	}

	return com.SuccessPayloadResponse(outputResponse)
}

func getAttachmentsByProcessReference(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	element := com.FPath.Path.PushBack("getAttachmentsByProcessReference")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		ProcessReference string `json:"processReference"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	var attachmentObjectsEntities []data.Entity
	response := data.QueryByIndex(&attachment.Attachment{}, &attachmentObjectsEntities, APIstub, []string{"ProcessId"}, request.ProcessReference)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	for i, attachmentObjectsEntity := range attachmentObjectsEntities {
		var attachmentVersionObjectEntities []data.Entity
		response := data.QueryByIndex(&version.Version{}, &attachmentVersionObjectEntities, APIstub, []string{"AttachmentId"}, attachmentObjectsEntity.GetId())
		if response.Status >= com.ERRORTHRESHOLD {
			return response
		}
		attachmentObjectsEntities[i].(*attachment.Attachment).NcVersions = data.EntitiesToOut(&version.Version{}, attachmentVersionObjectEntities)
	}

	return com.SuccessPayloadResponse(data.EntitiesToOut(&attachment.Attachment{}, attachmentObjectsEntities))
}

func getProcessElementHistory(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	element := com.FPath.Path.PushBack("getProcessElementHistory")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		ProcessReference  string `json:"processReference"`
		ElementNumber     int    `json:"elementNumber"`
		SinceResponseDate int64  `json:"sinceResponseDate"`
		ResponsesCount    int    `json:"responsesCount"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	process := process.Process{Id: request.ProcessReference}
	response := data.QueryById(&process, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	limitHistory := []proctemplate.Response{}
	if request.ResponsesCount < 0 {
		limitHistory = process.NcProcTemplate.Elements[request.ElementNumber+2].History
	} else {
		var startCount = 0
		for _, resp := range process.NcProcTemplate.Elements[request.ElementNumber+2].History {
			if startCount == 0 {
				if resp.ResponseSigned.DateTime < request.SinceResponseDate {
					if startCount < request.ResponsesCount {
						limitHistory = append(limitHistory, resp)
						startCount++
					} else {
						break
					}
				}
			} else {
				if startCount < request.ResponsesCount {
					limitHistory = append(limitHistory, resp)
					startCount++
				} else {
					break
				}
			}
		}
	}

	reverse := func(list []proctemplate.Response) []proctemplate.Response {
		for i := 0; i < len(list)/2; i++ {
			j := len(list) - i - 1
			list[i], list[j] = list[j], list[i]
		}
		return list
	}

	return com.SuccessPayloadResponse(reverse(limitHistory))
}

func addParticipantsArrToProcessObject(APIstub shim.ChaincodeStubInterface, processObject *process.Process) peer.Response {
	element := com.FPath.Path.PushBack("addParticipantsArrToProcessObject")
	defer com.FPath.Path.Remove(element)

	participants := []data.Entity{}
	for _, e := range processObject.NcProcTemplate.Elements {
		for ref, _ := range e.ParticipantsApproves {

			participantObject := participant.Participant{Id: ref}
			response := data.QueryById(&participantObject, APIstub)
			if response.Status >= com.ERRORTHRESHOLD {
				return response
			}

			participants = append(participants, &participantObject)
		}
	}
	processObject.Participants = data.EntitiesToOut(&participant.Participant{}, participants)

	return com.SuccessMessageResponse("All participants for process were gotten.")
}

func getAllProcesses(APIstub shim.ChaincodeStubInterface) peer.Response {
	element := com.FPath.Path.PushBack("getAllProcesses")
	defer com.FPath.Path.Remove(element)

	var processes []data.Entity
	response := data.QueryAll(&process.Process{}, &processes, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}
	executorRef, response := getExecutorFirstParticipantRef(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	myProcesses := []data.Entity{}

	for _, p := range processes {
		for _, e := range p.(*process.Process).NcProcTemplate.Elements {
			isFound := false
			for ref, _ := range e.ParticipantsApproves {
				if ref == executorRef {
					isFound = true
					break
				}
			}
			if isFound {
				myProcesses = append(myProcesses, p)
				break
			}
		}

		response := addParticipantsArrToProcessObject(APIstub, p.(*process.Process))
		if response.Status >= com.ERRORTHRESHOLD {
			return response
		}
	}

	return com.SuccessPayloadResponse(data.EntitiesToOut(&process.Process{}, myProcesses))
}

func createAttachment(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createAttachment")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		ProcessReference string `json:"processReference"`
		Name             string `json:"attachmentName"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	// Need to check have we got attachment with this name yet
	// Begin check name
	attachments := []data.Entity{}
	response := data.QueryByIndex(&attachment.Attachment{}, &attachments, APIstub, []string{"ProcessId"}, request.ProcessReference)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}
	for _, a := range attachments {
		if a.(*attachment.Attachment).Name == request.Name {
			return com.SuccessPayloadResponse(struct {
				Reference string `json:"reference"`
			}{Reference: a.(*attachment.Attachment).Id})
		}
	}
	// End check name

	attachmentObject := attachment.Attachment{Name: request.Name, ProcessId: request.ProcessReference}

	return data.Put(&attachmentObject, APIstub, simulate)
}

func createPart(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createPart")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		VersionReference string    `json:"versionReference"`
		Part             part.Part `json:"part"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	request.Part.VersionId = request.VersionReference

	return data.Put(&request.Part, APIstub, simulate)
}

func getLastAttachmentsVersions(APIstub shim.ChaincodeStubInterface, processReference string) ([]version.Version, peer.Response) {
	element := com.FPath.Path.PushBack("getLastAttachmentsVersions")
	defer com.FPath.Path.Remove(element)

	attachments := []data.Entity{}
	response := data.QueryByIndex(&attachment.Attachment{}, &attachments, APIstub, []string{"ProcessId"}, processReference)
	if response.Status >= com.ERRORTHRESHOLD {
		return nil, response
	}

	versions := []version.Version{}
	for _, a := range attachments {
		if a.(*attachment.Attachment).LastVersion != "" {
			version := version.Version{Id: a.(*attachment.Attachment).LastVersion}
			response := data.QueryById(&version, APIstub)
			if response.Status >= com.ERRORTHRESHOLD {
				return nil, response
			}
			versions = append(versions, version)
		}
	}

	return versions, com.SuccessMessageResponse("All last versions were gotten.")
}

func sendResponse(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("sendResponse")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		ProcessReference string                `json:"processReference"`
		ElementNumber    int                   `json:"elementNumber"`
		Response         proctemplate.Response `json:"response"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	ref, response := getExecutorFirstParticipantRef(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	time, response := funcs.GetTime(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	request.Response.ParticipantReference = ref
	request.Response.ResponseSigned.DateTime = time
	request.Response.ResponseSigned.ProcessReference = request.ProcessReference

	processObject := process.Process{Id: request.ProcessReference}

	response = data.QueryById(&processObject, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	// Start Set response in process
	elem := &processObject.NcProcTemplate.Elements[request.ElementNumber+2]

	// Check can executor work with process
	if _, ok := elem.ParticipantsApproves[ref]; !ok {
		return df_err.NoPermissionForWorkWithProcessError(ref)
	}

	// Set approves flags to approve map in process
	if request.Response.Type == proctemplate.ApproveType {
		// Set approve signatures to versions of files
		for _, vs := range request.Response.ResponseSigned.VersionSignatures {
			versionObject := version.Version{Id: vs.VersionReference}
			response := data.QueryById(&versionObject, APIstub)
			if response.Status >= com.ERRORTHRESHOLD {
				return response
			}
			if versionObject.Signatures == nil {
				com.DebugLogMsg("Test log for a wrong way! PEA")
				versionObject.Signatures = make(map[string]string)
			}
			versionObject.Signatures[ref] = vs.Signature
			response = data.EditAll(&versionObject, APIstub, simulate)
			if response.Status >= com.ERRORTHRESHOLD {
				return response
			}
		}

		// If need to approve all chat
		if request.Response.ResponseSigned.Approve == proctemplate.APPROVED {

			// Check are last versions of all files approved from executor
			versions, response := getLastAttachmentsVersions(APIstub, request.ProcessReference)
			if response.Status >= com.ERRORTHRESHOLD {
				return response
			}
			isAllApproved := true
			for _, v := range versions {
				if _, ok := v.Signatures[ref]; !ok {
					isAllApproved = false
					break
				}
			}

			// If all files are approved
			if isAllApproved {
				elem.ParticipantsApproves[ref] = request.Response.ResponseSigned.Approve
			}
		}

		// If need to disapprove all chat
		if request.Response.ResponseSigned.Approve == proctemplate.DISAPPROVED {
			elem.ParticipantsApproves[ref] = request.Response.ResponseSigned.Approve
		}
	}

	if request.Response.Type == proctemplate.FileType {
		for r, _ := range elem.ParticipantsApproves {
			elem.ParticipantsApproves[r] = proctemplate.NONE
		}
	}

	// Add response to process history
	elem.History = append([]proctemplate.Response{request.Response}, elem.History...)
	// End Set

	// Start Checking need we finish process or not
	goFinish := true
	for _, approve := range elem.ParticipantsApproves {
		if approve != proctemplate.APPROVED {
			goFinish = false
		}
	}

	if !goFinish {
		goFinish = true
		for _, approve := range elem.ParticipantsApproves {
			if approve != proctemplate.DISAPPROVED {
				goFinish = false
			}
		}
	}

	if goFinish {
		elem.Active = false
		processObject.NcProcTemplate.Elements[1].Active = true

		// Begin close event
		response = addParticipantsArrToProcessObject(APIstub, &processObject)
		if response.Status >= com.ERRORTHRESHOLD {
			return response
		}

		closeResponseEvect := struct {
			Response proctemplate.Response `json:"response"`
			Process  interface{}           `json:"process"`
		}{
			Response: request.Response,
			Process:  processObject.ToOut(),
		}

		eventJson, err := json.Marshal(closeResponseEvect)
		if err != nil {
			return com.MarshalError(err)
		}
		eventError := APIstub.SetEvent("close_process", eventJson)
		if err != nil {
			return com.EventError(eventError, "close_process")
		}
		// End close event
	} else {
		// Start Response Event
		eventJson, err := json.Marshal(&request.Response)
		if err != nil {
			return com.MarshalError(err)
		}
		eventError := APIstub.SetEvent("response", eventJson)
		if err != nil {
			return com.EventError(eventError, "response")
		}
		// End Event
	}
	// End Checking

	return data.EditAll(&processObject, APIstub, simulate)
}

func createAttachmentVersion(APIstub shim.ChaincodeStubInterface, args []string, simulate string) peer.Response {
	element := com.FPath.Path.PushBack("createAttachmentVersion")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		AttachmentReference string          `json:"attachmentReference"`
		Version             version.Version `json:"version"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	// Need to check hash in cuurent version of file and hash of new version
	// Begin check
	attachmentObject := attachment.Attachment{Id: request.AttachmentReference}
	response := data.QueryById(&attachmentObject, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	if attachmentObject.LastVersion != "" {
		currentVersion := version.Version{Id: attachmentObject.LastVersion}
		response = data.QueryById(&currentVersion, APIstub)
		if response.Status >= com.ERRORTHRESHOLD {
			return response
		}

		if currentVersion.Hash == request.Version.Hash {
			return com.SuccessMessageResponse("Hash of new version is equal to hash of current version. Version won't be edited.")
		}
	}
	// End check

	ref, response := getExecutorFirstParticipantRef(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	time, response := funcs.GetTime(APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	request.Version.AttachmentId = request.AttachmentReference
	request.Version.CreationTime = time
	request.Version.Creator = ref
	request.Version.Signatures = make(map[string]string)

	result := data.Put(&request.Version, APIstub, simulate)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	attachmentObject.LastVersion = request.Version.Id
	response = data.EditAll(&attachmentObject, APIstub, simulate)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return result
}

func getPartsByVersionReference(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	element := com.FPath.Path.PushBack("getPartsByVersionReference")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		VersionReference string `json:"versionReference"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	var partObjectsEntities []data.Entity
	response := data.QueryByIndex(&part.Part{}, &partObjectsEntities, APIstub, []string{"VersionId"}, request.VersionReference)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return com.SuccessPayloadResponse(data.EntitiesToOut(&attachment.Attachment{}, partObjectsEntities))
}

func getPartByReference(APIstub shim.ChaincodeStubInterface, args []string) peer.Response {
	element := com.FPath.Path.PushBack("getPartByReference")
	defer com.FPath.Path.Remove(element)

	if len(args) != 1 {
		return com.IncorrectInvokeNumberOfArgsError(args, 1)
	}

	type inputRequest struct {
		PartReference string `json:"partReference"`
	}

	var request inputRequest

	err := json.Unmarshal([]byte(args[0]), &request)
	if err != nil {
		return com.UnmarshalError(err, args[0])
	}

	var partObject = part.Part{Id: request.PartReference}
	response := data.QueryById(&partObject, APIstub)
	if response.Status >= com.ERRORTHRESHOLD {
		return response
	}

	return com.SuccessPayloadResponse(partObject.ToSoloOut())
}
