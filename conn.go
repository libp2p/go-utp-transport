package utp

// This file manage the multiaddr friendly conn object

import (
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

type utpMaConn struct {
	net.Conn

	laddr ma.Multiaddr
	raddr ma.Multiaddr
}

// Wrap a net.Conn with laddr and raddr known
func newMaConnWithAddrs(c net.Conn, laddr ma.Multiaddr, raddr ma.Multiaddr) utpMaConn {
	return utpMaConn{
		Conn:  c,
		laddr: laddr,
		raddr: raddr,
	}
}

// Create with raddr known (for dial out)
func newMaConnWithRaddr(c net.Conn, raddr ma.Multiaddr) (utpMaConn, error) {
	// Convert the conn to a Multiaddr, this will only do ip and udp.
	mabase, err := manet.FromNetAddr(c.LocalAddr())
	if err != nil {
		return utpMaConn{}, err
	}

	return newMaConnWithAddrs(
		c,
		mabase.Encapsulate(emptyUtpMa),
		raddr,
	), nil
}

// Create with laddr known (for Accepting in)
func newMaConnWithLaddr(c net.Conn, laddr ma.Multiaddr) (utpMaConn, error) {
	// Convert the conn to a Multiaddr, this will only do ip and udp.
	mabase, err := manet.FromNetAddr(c.RemoteAddr())
	if err != nil {
		return utpMaConn{}, err
	}

	return newMaConnWithAddrs(
		c,
		laddr,
		mabase.Encapsulate(emptyUtpMa),
	), nil
}

func (c utpMaConn) LocalMultiaddr() ma.Multiaddr {
	return c.laddr
}

func (c utpMaConn) RemoteMultiaddr() ma.Multiaddr {
	return c.raddr
}
