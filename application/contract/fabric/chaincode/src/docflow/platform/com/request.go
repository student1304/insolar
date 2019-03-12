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

type Request struct {
	Channels []string `json:"channels>channel"`

	EntityName string `json:"entity_name"`
	Type       string `json:"type"`
	TypeAttr   string `json:"type_attr"`

	Args map[string]string `json:"args>arg"`

	// Edit Fields //////////////////////////////
	FieldPaths  []FieldPath `json:"field_paths>field_path"`
	FieldValues []string    `json:"field_values>field_value"`
	/////////////////////////////////////////////

	Simulate string `json:"simulate"`
}

type FieldPath struct {
	FieldPath []string `json:"field"`
}
