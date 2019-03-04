package fabric

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Fabric struct {
	foundation.BaseContract
	transactionId int
	State         map[string][]byte
	InsStub       shim.InsStub
}

func New(name string, cc shim.Chaincode) (*Fabric, error) {
	state := make(map[string][]byte)
	return &Fabric{
		transactionId: 0,
		State:         state,
		InsStub:       *shim.NewInsStub(name, cc, &state),
	}, nil
}

var INSATTR_Call_API = true

// Call method for authorized calls
func (f *Fabric) Call(rootDomain core.RecordRef, method string, params []byte, seed []byte, sign []byte) (interface{}, error) {

	switch method {
	case "Init":
		response := f.InsStub.MockInit(string(f.transactionId), [][]byte{params})
		if response.Status >= 400 {
			return response.Payload, fmt.Errorf("[ Call ]: init error: %s", response.Message)
		} else {
			return response.Payload, nil
		}
	case "Invoke":
		response := f.InsStub.MockInvoke(string(f.transactionId), [][]byte{params})
		if response.Status >= 400 {
			return response.Payload, fmt.Errorf("[ Call ]: invoke error: %s", response.Message)
		} else {
			return response.Payload, nil
		}
	}
	return nil, &foundation.Error{S: "Unknown method"}
}
