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

package participant

import "github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/record"

type Participant struct {
	RelationTxId string `json:"relationTxId"`
	Key          string `json:"key"`
	Id           string `json:"id"`

	PersonName       string `json:"personName"`
	OrganizationName string `json:"organizationName"`
	Position         string `json:"position"`
	MspId            string `json:"mspId"`
	PublicKey        string `json:"publicKey"`
	Email            string `json:"email"`
}

type ParticipantOut struct {
	Email            string `json:"email"`
	MspId            string `json:"mspId"`
	OrganizationName string `json:"organizationName"`
	PersonName       string `json:"personName"`
	Position         string `json:"position"`
	PublicKey        string `json:"publicKey"`
}

func (participant Participant) ToOut() interface{} {
	participantOut := ParticipantOut{
		PersonName:       participant.PersonName,
		OrganizationName: participant.OrganizationName,
		Position:         participant.Position,
		PublicKey:        participant.PublicKey,
		Email:            participant.Email,
		MspId:            participant.MspId,
	}

	record := record.Record{
		ParentReference: "",
		Reference:       participant.Id,
		Object:          participantOut,
	}

	return record
}
