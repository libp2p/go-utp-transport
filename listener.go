package utp

// This file manage to multiaddr friendly listener.

import (
	"net"

	"github.com/anacrolix/utp"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type utpMaListener struct {
	// The utp listener
	list *utp.Socket

	// The local multiaddr address
	laddr ma.Multiaddr
}

func newListener(l *utp.Socket) (utpMaListener, error) {
	// Recreating the laddr is to get the real if listening on :0.
	mabase, err := manet.FromNetAddr(l.Addr())
	if err != nil {
		return utpMaListener{}, err
	}
	return utpMaListener{
		list: l,
		// Adding "/utp".
		laddr: mabase.Encapsulate(emptyUtpMa),
	}, nil
}

func (l utpMaListener) Accept() (manet.Conn, error) {
	// Accepting the conn
	utpconn, err := l.list.Accept()
	if err != nil {
		return nil, err
	}

	// Wrapping it into a multiaddr
	return newMaConnWithLaddr(utpconn, l.laddr)
}

func (l utpMaListener) Close() error {
	return l.list.Close()
}

func (l utpMaListener) Multiaddr() ma.Multiaddr {
	return l.laddr
}

func (l utpMaListener) Addr() net.Addr {
	return l.list.Addr()
}
