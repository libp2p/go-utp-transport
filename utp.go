package utp

// This file manage the network part.

import (
	"context"
	"net"

	"github.com/anacrolix/utp"
	"github.com/libp2p/go-libp2p-core/peer"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	tptu "github.com/libp2p/go-libp2p-transport-upgrader"
	ma "github.com/multiformats/go-multiaddr"
	mafmt "github.com/multiformats/go-multiaddr-fmt"
	manet "github.com/multiformats/go-multiaddr-net"
)

var emptyUtpMa, _ = ma.NewMultiaddr("/utp")

var _ tpt.Transport = (*UtpTransport)(nil)

type UtpTransport struct {
	// The upgrader upgrade connections from raw utp to full saffed multiplexed ones.
	// I'm not sure if the default encryption and multiplexer are capable to deal well with utp, need test.
	Upgrader *tptu.Upgrader

	// The socket is the utp listener connection, its reused to dial.
	socket *utp.Socket
}

func NewUTPTransport(u *tptu.Upgrader) *UtpTransport {
	return &UtpTransport{Upgrader: u, socket: nil}
}

func (t *UtpTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (tpt.CapableConn, error) {
	// Converting multiaddr to net addr
	network, addr, err := manet.DialArgs(raddr.Decapsulate(emptyUtpMa))
	if err != nil {
		return nil, err
	}
	var utpconn net.Conn
	// Doing the utp connection
	// Check for an open listener
	if t.socket == nil {
		// If not creating one with a new filedescriptor
		utpconn, err = utp.DialContext(ctx, addr)
	} else {
		// If using the open one
		utpconn, err = t.socket.DialContext(ctx, network, addr)
	}
	if err != nil {
		return nil, err
	}
	// Transforming it to an multiaddr conn
	maconn, err := manet.WrapNetConn(utpconn)
	if err != nil {
		return nil, err
	}
	// Upgrading it to an safe mutliplexed one and return (the conn and the err).
	return t.Upgrader.UpgradeOutbound(ctx, t, maconn, p)
}

func (t *UtpTransport) CanDial(addr ma.Multiaddr) bool {
	return mafmt.UTP.Matches(addr)
}

func (t *UtpTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {
	// Converting multiaddr to net addr
	network, addr, err := manet.DialArgs(laddr.Decapsulate(emptyUtpMa))
	if err != nil {
		return nil, err
	}
	// Creating the socket
	utpsock, err := utp.NewSocket(network, addr)
	if err != nil {
		return nil, err
	}
	// Wrapping it into an multiaddr one
	malist, err := manet.WrapNetListener(utpsock)
	if err != nil {
		return nil, err
	}
	// Adding it for reusing in dial
	t.socket = utpsock
	// Upgrading the listener to an safe and multiplexed one and return.
	return t.Upgrader.UpgradeListener(t, malist), nil
}

func (t *UtpTransport) Protocols() []int {
	return []int{ma.P_UTP}
}

func (t *UtpTransport) Proxy() bool {
	return false
}
