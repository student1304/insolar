// /
//    Copyright 2019 Insolar
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
// /

package platformpolicy

import (
	"crypto/rand"
	"testing"

	"github.com/insolar/x-crypto/elliptic"
	"github.com/stretchr/testify/assert"
)

var p256KeyProcessor = &keyProcessor{
	curve: elliptic.P256(),
}

var secp256k1KeyProcessor = &keyProcessor{
	curve: elliptic.Secp256k1(),
}

func BenchmarkP256Sign1000(b *testing.B) {
	benchmarkSign(b, p256KeyProcessor, 1000)
}

func BenchmarkP256k1Sign1000(b *testing.B) {
	benchmarkSign(b, secp256k1KeyProcessor, 1000)
}

func BenchmarkP256Sign10000(b *testing.B) {
	benchmarkSign(b, p256KeyProcessor, 10000)
}

func BenchmarkP256k1Sign10000(b *testing.B) {
	benchmarkSign(b, secp256k1KeyProcessor, 10000)
}

func benchmarkSign(b *testing.B, processor *keyProcessor, length int) {
	sk, err := processor.GeneratePrivateKey()
	assert.NoError(b, err)

	cs := NewPlatformCryptographyScheme()
	payload := make([]byte, length)
	_, err = rand.Read(payload)
	assert.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, err = cs.Signer(sk).Sign(payload)
	}
	assert.NoError(b, err)
}
