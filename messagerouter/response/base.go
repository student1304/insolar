/*
 *    Copyright 2018 Insolar
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

// Package response represents responses to messages of the messagerouter
package response

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type baseResponse struct {
}

// Type is a enum type of response
type Type byte

const (
	WrongResponseType  = Type(iota)
	CommonResponseType // CommonResponseType - two binary fields: data and results
)

func getEmptyResponse(t Type) (core.Response, error) {
	switch t {
	case WrongResponseType:
		return nil, errors.New("no empty response for 'wrong' response")
	case CommonResponseType:
		return &CommonResponse{}, nil
	default:
		return nil, errors.Errorf("unimplemented response type: '%d'", t)
	}
}

func serialize(m core.Response, t Type) (io.Reader, error) {
	buff := &bytes.Buffer{}
	_, err := buff.Write([]byte{byte(t)})
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(buff)
	err = enc.Encode(m)
	return buff, err
}

// Deserialize returns a response
func Deserialize(buff io.Reader) (core.Response, error) {
	b := make([]byte, 1)
	_, err := buff.Read(b)
	if err != nil {
		return nil, errors.New("too short input to deserialize a message response")
	}

	m, err := getEmptyResponse(Type(b[0]))
	if err != nil {
		return nil, err
	}
	enc := gob.NewDecoder(buff)
	err = enc.Decode(m)
	return m, err
}

func init() {
	gob.Register(&CommonResponse{})
}
