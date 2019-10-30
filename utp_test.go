package utp

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"

	"github.com/libp2p/go-libp2p-core/sec/insecure"
	mplex "github.com/libp2p/go-libp2p-mplex"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	ma "github.com/multiformats/go-multiaddr"

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
	ttpt.SubtestTransport(t, ta, tb, zero, ia.LocalPeer())
}

// The utp lib was not dealing well with overload so this test is specificaly
// here to ensure it does, even if that not realy a bug of the transport we need
// to be sure it can handle the thousands connection IPFS or any other p2p app do.
func TestUtpTransportOverload(t *testing.T) {
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

	zero, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/0/utp")
	if err != nil {
		t.Fatal(err)
	}
	ttpt.SubtestStress(t, ta, tb, zero, ia.LocalPeer(), ttpt.Options{
		// Don't put too much to avoid running out of file descriptor, the kernel
		// doesn't like listening on 1000 ports.
		ConnNum:   250,
		StreamNum: 2,
		MsgNum:    50,
		MsgMax:    100,
		MsgMin:    100,
	})
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
