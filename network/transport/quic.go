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

package transport

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"net"

	"github.com/insolar/insolar/log"
	"github.com/insolar/insolar/network/transport/relay"
	"github.com/lucas-clemente/quic-go"
)

type quicTransport struct {
	baseTransport
	l quic.Listener
}

func newQuicTransport(conn net.PacketConn, proxy relay.Proxy, publicAddress string) (*quicTransport, error) {
	listener, err := quic.Listen(conn, generateTLSConfig(), nil)
	if err != nil {
		return nil, err
	}

	transport := &quicTransport{
		baseTransport: newBaseTransport(proxy, publicAddress),
		l:             listener,
	}

	transport.sendFunc = transport.send
	return transport, nil
}

func (q *quicTransport) send(recvAddress string, data []byte) error {
	ctx := context.Background()
	session, err := quic.DialAddrContext(ctx, recvAddress, &tls.Config{InsecureSkipVerify: true}, nil)
	if err != nil {
		return err
	}
	//defer session.Close()

	log.Infof("connected to: %s", session.RemoteAddr().String())

	stream, err := session.OpenStreamSync()
	if err != nil {
		return err
	}
	//defer stream.Close()

	_, err = stream.Write(data)
	if err != nil {
		return err
	}

	return err
}

// Start starts networking.
func (q *quicTransport) Start(ctx context.Context) error {
	log.Info("Start QUIC transport")
	for {
		session, err := q.l.Accept()
		if err != nil {
			<-q.disconnectFinished
			return err
		}

		log.Infof("accept from: %s", session.RemoteAddr().String())
		go q.handleAcceptedConnection(session)
	}
}

// Stop stops networking.
func (q *quicTransport) Stop() {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	log.Info("Stop QUIC transport")
	q.prepareDisconnect()

	err := q.l.Close()
	if err != nil {
		log.Errorln("Failed to close socket:", err.Error())
	}
}

func (q *quicTransport) handleAcceptedConnection(session quic.Session) {
	defer session.Close()

	stream, err := session.AcceptStream()

	msg, err := q.serializer.DeserializePacket(stream)
	if err != nil {
		log.Error("[ handleAcceptedConnection ] ", err)
		return
	}

	q.handlePacket(msg)

}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{Certificates: []tls.Certificate{tlsCert}}
}
