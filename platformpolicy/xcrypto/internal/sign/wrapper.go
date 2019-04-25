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
	"crypto/rand"
	"encoding/asn1"
	"math/big"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/insolar/x-crypto/ecdsa"
)

type ecdsaSignerWrapper struct {
	privateKey *ecdsa.PrivateKey
	hasher     insolar.Hasher
}

func (sw *ecdsaSignerWrapper) Sign(data []byte) (*insolar.Signature, error) {
	hash := sw.hasher.Hash(data)

	r, s, err := ecdsa.Sign(rand.Reader, sw.privateKey, hash)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	ecdsaSignature := SerializeTwoBigInt(r, s)
	if err != nil {
		return nil, errors.Wrap(err, "[ Sign ] could't sign data")
	}

	signature := insolar.SignatureFromBytes(ecdsaSignature)
	return &signature, nil
}

type ecdsaVerifyWrapper struct {
	publicKey *ecdsa.PublicKey
	hasher    insolar.Hasher
}

func (sw *ecdsaVerifyWrapper) Verify(signType string, signature insolar.Signature, data []byte) (bool, error) {
	if signature.Bytes() == nil {
		return false, errors.Errorf("[ Verify ] signature bytes = nil")
	}

	var r, s *big.Int
	var err error
	switch signType {
	case "DER":
		r, s, err = PointsFromDER(signature.Bytes())
		if err != nil {
			log.Error(err)
			return false, errors.Wrap(err, "[ Verify ] could't get points from DER")
		}
	default:
		r, s, err = DeserializeTwoBigInt(signature.Bytes())
		if err != nil {
			log.Error(err)
			return false, errors.Wrap(err, "[ Verify ] could't deserialize signature")
		}
	}

	hash := sw.hasher.Hash(data)
	return ecdsa.Verify(sw.publicKey, hash, r, s), nil
}

// Get the X and Y points from a DER encoded signature
// Sometimes demarshalling using Golang's DEC to struct unmarshalling fails; this extracts R and S from the bytes
// manually to prevent crashing.
// This should NOT be a hex encoded byte array
func PointsFromDER(der []byte) (R, S *big.Int, err error) {
	R, S = &big.Int{}, &big.Int{}

	data := asn1.RawValue{}
	if _, err := asn1.Unmarshal(der, &data); err != nil {
		return &big.Int{}, &big.Int{}, errors.Wrap(err, "[ PointsFromDER ] Unmarshal der error: ")
	}

	// The format of our DER string is 0x02 + rlen + r + 0x02 + slen + s
	rLen := data.Bytes[1] // The entire length of R + offset of 2 for 0x02 and rlen
	r := data.Bytes[2 : rLen+2]
	// Ignore the next 0x02 and slen bytes and just take the start of S to the end of the byte array
	s := data.Bytes[rLen+4:]

	R.SetBytes(r)
	S.SetBytes(s)

	return
}
