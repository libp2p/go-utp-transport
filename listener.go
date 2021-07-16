package utp

// This file manage to multiaddr friendly listener.

import (
	"net"

	"github.com/anacrolix/go-libutp"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type utpMaListener struct {
	// The utp listener
	list *utp.Socket

	// The parent transport, used to remove the listener reference for filedescriptor reusing at closing.
	t *UtpTransport

	// The local multiaddr address
	laddr ma.Multiaddr
}

func newListener(l *utp.Socket, t *UtpTransport) (utpMaListener, error) {
	// Recreating the laddr is to get the real port if listening on :0.
	mabase, err := manet.FromNetAddr(l.Addr())
	if err != nil {
		return utpMaListener{}, err
	}
	return utpMaListener{
		list: l,
		t:    t,
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
	// Checking for the ip version
	if l.laddr.Protocols()[0].Code == ma.P_IP4 {
		// Locking the file descriptor locker
		l.t.sl4.Lock()
		// Checking if the socket used for file descriptor reuse is the same as the one we are closing
		if l.t.socket4 == l.list {
			// If removing the reference to avoid reusing a closed socket
			l.t.socket4 = nil
		}
		// Unlocking
		l.t.sl4.Unlock()
	} else {
		// Locking the file descriptor locker
		l.t.sl6.Lock()
		// Checking if the socket used for file descriptor reuse is the same as the one we are closing
		if l.t.socket6 == l.list {
			// If removing the reference to avoid reusing a closed socket
			l.t.socket6 = nil
		}
		// Unlocking
		l.t.sl6.Unlock()
	}
	// Closing
	return l.list.Close()
}

func (l utpMaListener) Multiaddr() ma.Multiaddr {
	return l.laddr
}

func (l utpMaListener) Addr() net.Addr {
	return l.list.Addr()
}
