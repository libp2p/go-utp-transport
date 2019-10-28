package utp

import (
	"testing"

	ma "github.com/multiformats/go-multiaddr"
)

func TestCanDial(t *testing.T) {
	addrUtp, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/5555/utp")
	if err != nil {
		t.Fatal(err)
	}
	addrTCP, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5555")
	if err != nil {
		t.Fatal(err)
	}
	addrTCPUtp, err := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/5555/utp")
	if err != nil {
		t.Fatal(err)
	}

	d := &UtpTransport{}
	matchTrue := d.CanDial(addrUtp)
	matchFalse := d.CanDial(addrTCP)
	matchFalseTcp := d.CanDial(addrTCPUtp)

	if !matchTrue {
		t.Fatal("expected to match udp+utp maddr, but did not")
	}
	if matchFalse {
		t.Fatal("expected to not match tcp maddr, but did")
	}
	if matchFalseTcp {
		t.Fatal("expected to not match tcp+utp maddr, but did")
	}
}
