package utp

import (
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"testing"

	"github.com/libp2p/go-libp2p-core/sec/insecure"
	mplex "github.com/libp2p/go-libp2p-mplex"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"

	ttpt "github.com/libp2p/go-libp2p-testing/suites/transport"

	ma "github.com/multiformats/go-multiaddr"
)

func TestUtpTransport(t *testing.T) {
	ia := makeInsecureTransport(t)
	ib := makeInsecureTransport(t)

	ta := NewUTPTransport(&tptu.Upgrader{
		Secure: ia,
		Muxer:  new(mplex.Transport),
	})
	tb := NewUTPTransport(&tptu.Upgrader{
		Secure: ib,
		Muxer:  new(mplex.Transport),
	})

	zero := "/ip4/127.0.0.1/udp/0/utp"
	ttpt.SubtestTransport(t, ta, tb, zero, ia.LocalPeer())
}

func TestUtpTransportCantListenUtp(t *testing.T) {
	utpa, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/0")
	if err != nil {
		t.Fatal(err)
	}

	tpt := NewUTPTransport(&tptu.Upgrader{
		Secure: makeInsecureTransport(t),
		Muxer:  new(mplex.Transport),
	})

	_, err = tpt.Listen(utpa)
	if err == nil {
		t.Fatal("shouldnt be able to listen on utp addr with utp transport")
	}
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
