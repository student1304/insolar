/*
 *    Copyright 2019 Insolar Technologies
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

package df_templates

import (
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/docflow/df_data/data/proctemplate"
)

var (
	ChatTemplate = proctemplate.ProcTemplate{
		RelationTxId: "11111",
		NcProcTemplate: proctemplate.NcProcTemplate{
			Name: "Chat",
			Elements: []proctemplate.Element{
				{
					Name:                 "Start",
					ParticipantsApproves: make(map[string]proctemplate.Approve),
					Active:               true,
					History:              []proctemplate.Response{},
				},
				{
					Name:                 "Finish",
					ParticipantsApproves: make(map[string]proctemplate.Approve),
					Active:               false,
					History:              []proctemplate.Response{},
				},
				{
					Name:                 "ChatActivity",
					ParticipantsApproves: make(map[string]proctemplate.Approve),
					Active:               false,
					History:              []proctemplate.Response{},
				},
			},
		},
	}
)
