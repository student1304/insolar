//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package helloworld

import (
	"fmt"

	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

// HelloWorld contract
type HelloWorld struct {
	foundation.BaseContract
	Greeted int
}

var INSATTR_Greet_API = true

func (hw *HelloWorld) Call()

// Greet greats the caller
func (hw *HelloWorld) Greet(name string) (interface{}, error) {
	hw.Greeted++
	return fmt.Sprintf("Hello %s' world", name), nil
}

// New returns a new empty contract
func New() (*HelloWorld, error) {
	return &HelloWorld{
		Greeted: 0,
	}, nil
}
