package fabric

import (
	"fmt"
	"github.com/insolar/insolar/application/contract/fabric/chaincode/src/docflow/platform/sc"
	"github.com/insolar/insolar/application/contract/fabric/insstub"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type Fabric struct {
	foundation.BaseContract
	transactionId int
	State         map[string][]byte
	InsStub       insstub.InsStub
}

func New() (*Fabric, error) {
	state := make(map[string][]byte)
	return &Fabric{
		transactionId: 0,
		State:         state,
		InsStub:       *insstub.NewInsStub("insFabric", new(sc.SmartContract), &state),
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
