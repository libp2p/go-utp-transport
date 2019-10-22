package utp

import (
	"net"

	"github.com/libp2p/go-libp2p-core/transport"
	ma "github.com/multiformats/go-multiaddr"
)

type UtpListener struct {
}

func (l *UtpListener) Accept() (transport.CapableConn, error) {
	return nil, nil
}

func (l *UtpListener) Close() error {
	return nil
}

func (l *UtpListener) Addr() net.Addr {
	return nil
}

func (l *UtpListener) Multiaddr() ma.Multiaddr {
	return nil
}
