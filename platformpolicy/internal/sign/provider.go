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

package sign

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy/internal/hash"
)

type ecdsaProvider struct {
	HashProvider hash.AlgorithmProvider `inject:""`
}

func NewECDSAProvider() AlgorithmProvider {
	return &ecdsaProvider{}
}

func (p *ecdsaProvider) Sign(privateKey crypto.PrivateKey) insolar.Signer {
	return &ecdsaSignerWrapper{
		privateKey: MustConvertPrivateKeyToEcdsa(privateKey),
		hasher:     p.HashProvider.Hash512bits(),
	}
}

func (p *ecdsaProvider) Verify(publicKey crypto.PublicKey) insolar.Verifier {
	return &ecdsaVerifyWrapper{
		publicKey: MustConvertPublicKeyToEcdsa(publicKey),
		hasher:    p.HashProvider.Hash512bits(),
	}
}