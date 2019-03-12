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

const (
	// Name of entity (for logs)
	ENTITY_NAME = "ParticipantReference"

	// JSON tag name
	JSON_TAG = "participants"

	// Object type names (for storage)
	KEY = "PARTICIPANT"
)

func (participant Participant) CreateValidation() bool {
	return true
}

func (participant Participant) ChangeValidation(newParticipantInterface interface{}) bool {
	_ = newParticipantInterface.(*Participant)
	return true
}

func (participant Participant) GetKeyObjectType() string {
	return KEY
}

func (participant Participant) GetIndexes() [][]string {
	return [][]string{
		{"MspId"},
	}
}

func (participant Participant) GetEntityName() string {
	return ENTITY_NAME
}
func (participant Participant) GetTagName() string {
	return JSON_TAG
}
func (participant *Participant) SetTxId(relationTxId string) {
	participant.RelationTxId = relationTxId
}
func (participant *Participant) GetTxId() string {
	return participant.RelationTxId
}
func (participant *Participant) SetId(id string) {
	participant.Id = id
}
func (participant Participant) GetId() string {
	return participant.Id
}
func (participant *Participant) SetKey(key string) {
	participant.Key = key
}
func (participant Participant) GetKey() string {
	return participant.Key
}
