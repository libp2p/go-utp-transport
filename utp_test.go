package utp

import (
	"testing"
	"time"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec/insecure"
	mplex "github.com/libp2p/go-libp2p-mplex"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"

	ttpt "github.com/libp2p/go-libp2p-testing/suites/transport"
)

func TestUtpTransport(t *testing.T) {
	ia := makeInsecureTransport(t)
	ib := makeInsecureTransport(t)
	ta := NewUtpTransport(&tptu.Upgrader{
		Secure: ia,
		Muxer:  new(mplex.Transport),
	})
	tb := NewUtpTransport(&tptu.Upgrader{
		Secure: ib,
		Muxer:  new(mplex.Transport),
	})

	zero := "/ip4/127.0.0.1/udp/0/utp"
	ttpt.SubtestTransportThrottled(t, ta, tb, zero, ia.LocalPeer(), 15 * time.Millisecond)
}

func makeInsecureTransport(t *testing.T) *insecure.Transport {
	priv, pub, err := crypto.GenerateKeyPair(crypto.Ed25519, 256)
	if err != nil {
		t.Fatal(err)
	}
	id, err := peer.IDFromPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return insecure.NewWithIdentity(id, priv)
}
