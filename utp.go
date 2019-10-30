package utp

// This file manage the network part.

import (
	"context"
	"net"
	"sync"

	"github.com/anacrolix/go-libutp"
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

	// Socket Locker (ip 4 or 6)
	// The RWMutex is to avoid race issue.
	// Multiple dial can occurs that why its RWMutex.
	sl4 sync.RWMutex
	sl6 sync.RWMutex

	// This is the socket to reuse for dialing.
	// Only socket with laddr passing manet.IPUnspecified() == true can be putted here to avoid to select a working socket beetween all avaible socket.
	socket4 *utp.Socket
	socket6 *utp.Socket
}

func NewUtpTransport(u *tptu.Upgrader) *UtpTransport {
	return &UtpTransport{Upgrader: u}
}

func (t *UtpTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (tpt.CapableConn, error) {
	// Converting multiaddr to net addr
	network, addr, err := manet.DialArgs(raddr.Decapsulate(emptyUtpMa))
	if err != nil {
		return nil, err
	}

	var utpconn net.Conn
	// Doing the utp connection
	// Check for the network
	if raddr.Protocols()[0].Code == ma.P_IP4 {
		// Lock in read mode to avoid race with listener
		t.sl4.RLock()
		// Check for an open listener
		if t.socket4 == nil {
			// Unlock as read a lock as write, that shouldn't happend more than
			// twice per execution so that not a really high cost
			t.sl4.RUnlock()
			t.sl4.Lock()
			// Cheking again because it was unlocked for a moment
			if t.socket4 == nil {
				// If not creating a new listener not for listening, only dial
				newList, err := utp.NewSocket("udp4", "0.0.0.0:0")
				if err != nil {
					t.sl4.Unlock()
					return nil, err
				}
				// If there was no error setting a new socket
				t.socket4 = newList
			}
			// Unlock and relock as read
			t.sl4.Unlock()
			t.sl4.RLock()
		}
		// If using the open one
		utpconn, err = t.socket4.DialContext(ctx, network, addr)
		// We have our socket, unlocking
		t.sl4.RUnlock()
	} else {
		// Lock in read mode to avoid race with listener
		t.sl6.RLock()
		// Check for an open listener
		if t.socket6 == nil {
			// Unlock as read a lock as write, that shouldn't happend more than
			// twice per execution so that not a really high cost
			t.sl6.RUnlock()
			t.sl6.Lock()
			// Cheking again because it was unlocked for a moment
			if t.socket6 == nil {
				// If not creating a new listener not for listening, only dial
				newList, err := utp.NewSocket("udp6", "[::]:0")
				if err != nil {
					t.sl6.Unlock()
					return nil, err
				}
				// If there was no error setting a new socket
				t.socket6 = newList
			}
			// Unlock and relock as read
			t.sl6.Unlock()
			t.sl6.RLock()
		}
		// If using the open one
		utpconn, err = t.socket6.DialContext(ctx, network, addr)
		// We have our socket, unlocking
		t.sl6.RUnlock()
	}
	if err != nil {
		return nil, err
	}

	// Transforming it to an multiaddr conn
	maconn, err := newMaConnWithRaddr(utpconn, raddr)
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
	malist, err := newListener(utpsock, t)
	if err != nil {
		return nil, err
	}
	// Adding it for reusing in dial
	// Check if that IPUnspecified
	if manet.IsIPUnspecified(laddr) {
		// Check for ip version
		if laddr.Protocols()[0].Code == ma.P_IP4 {
			// Lock to avoid race with other listener or dialer
			t.sl4.Lock()
			// Cheking for an already setuped reuse
			if t.socket4 == nil {
				// Setting up reuse
				t.socket4 = utpsock
			}
			t.sl4.Unlock()
		} else {
			// Lock to avoid race with other listener or dialer
			t.sl6.Lock()
			// Cheking for an already setuped reuse
			if t.socket6 == nil {
				// Setting up reuse
				t.socket6 = utpsock
			}
			t.sl6.Unlock()
		}
	}
	// Upgrading the listener to an safe and multiplexed one and return.
	return t.Upgrader.UpgradeListener(t, malist), nil
}

func (t *UtpTransport) Protocols() []int {
	return []int{ma.P_UTP}
}

func (t *UtpTransport) Proxy() bool {
	return false
}
