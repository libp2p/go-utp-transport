package utp

import (
	"context"

	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
	mafmt "github.com/multiformats/go-multiaddr-fmt"
)

type UtpTransport struct {
}

func (t *UtpTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (transport.CapableConn, error) {
	return nil, nil
}

func (t *UtpTransport) CanDial(addr ma.Multiaddr) bool {
	return mafmt.UTP.Matches(addr)
}

func (t *UtpTransport) Listen(laddr ma.Multiaddr) (UtpListener, error) {
	return UtpListener{}, nil
}

func (t *UtpTransport) Protocols() []int {
	return []int{ma.P_UTP}
}

func (t *UtpTransport) Proxy() bool {
	return false
}
