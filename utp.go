package transport

import (
	"context"
	"fmt"
	"net"
	"sync"

	utp "github.com/anacrolix/utp"
	mafmt "gx/ipfs/QmQkdkvXE4oKXAcLZK5d7Zc6xvyukQc8WVjX7QvxDJ7hJj/mafmt"
	manet "gx/ipfs/QmT6Cp31887FpAc25z25YHgpFJohZedrYLWPPspRtj1Brp/go-multiaddr-net"
	ma "gx/ipfs/QmUAQaWbKxGCUTuoQVvvicbQNZ9APF5pDGWyAZSe93AtKH/go-multiaddr"
	tpt "gx/ipfs/QmWMia2fBVBesMerbtApQY7Tj2sgTaziveBACfCRUcv45f/go-libp2p-transport"
)

var errIncorrectNetAddr = fmt.Errorf("incorrect network addr conversion")

var utpAddrSpec = &manet.NetCodec{
	ProtocolName:     "utp",
	NetAddrNetworks:  []string{"utp", "utp4", "utp6"},
	ParseNetAddr:     parseUtpNetAddr,
	ConvertMultiaddr: parseUtpMaddr,
}

func init() {
	manet.RegisterNetCodec(utpAddrSpec)
}

type UtpTransport struct {
	sockLock sync.Mutex
	sockets  map[string]*UtpSocket
}

func NewUtpTransport() *UtpTransport {
	return &UtpTransport{
		sockets: make(map[string]*UtpSocket),
	}
}

func (d *UtpTransport) Matches(a ma.Multiaddr) bool {
	return mafmt.UTP.Matches(a)
}

type UtpSocket struct {
	s         *utp.Socket
	laddr     ma.Multiaddr
	transport tpt.Transport
}

func (t *UtpTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {
	t.sockLock.Lock()
	defer t.sockLock.Unlock()
	s, ok := t.sockets[laddr.String()]
	if ok {
		return s, nil
	}

	ns, err := t.newConn(laddr)
	if err != nil {
		return nil, err
	}

	t.sockets[laddr.String()] = ns
	return ns, nil
}

func (t *UtpTransport) Dialer(laddr ma.Multiaddr, opts ...tpt.DialOpt) (tpt.Dialer, error) {
	t.sockLock.Lock()
	defer t.sockLock.Unlock()
	s, ok := t.sockets[laddr.String()]
	if ok {
		return s, nil
	}

	ns, err := t.newConn(laddr, opts...)
	if err != nil {
		return nil, err
	}

	t.sockets[laddr.String()] = ns
	return ns, nil
}

func (t *UtpTransport) newConn(addr ma.Multiaddr, opts ...tpt.DialOpt) (*UtpSocket, error) {
	network, netaddr, err := manet.DialArgs(addr)
	if err != nil {
		return nil, err
	}

	s, err := utp.NewSocket("udp"+network[3:], netaddr)
	if err != nil {
		return nil, err
	}

	laddr, err := manet.FromNetAddr(s.LocalAddr())
	if err != nil {
		return nil, err
	}

	return &UtpSocket{
		s:         s,
		laddr:     laddr,
		transport: t,
	}, nil
}

func (s *UtpSocket) Dial(raddr ma.Multiaddr) (tpt.Conn, error) {
	return s.DialContext(context.Background(), raddr)
}

func (s *UtpSocket) DialContext(ctx context.Context, raddr ma.Multiaddr) (tpt.Conn, error) {
	_, addr, err := manet.DialArgs(raddr)
	if err != nil {
		return nil, err
	}

	// TODO: update utp lib
	con, err := s.s.Dial(addr)
	if err != nil {
		return nil, err
	}

	mnc, err := manet.WrapNetConn(con)
	if err != nil {
		return nil, err
	}

	return &tpt.ConnWrap{
		Conn: mnc,
		Tpt:  s.transport,
	}, nil
}

func (s *UtpSocket) Accept() (tpt.Conn, error) {
	c, err := s.s.Accept()
	if err != nil {
		return nil, err
	}

	mnc, err := manet.WrapNetConn(c)
	if err != nil {
		return nil, err
	}

	return &tpt.ConnWrap{
		Conn: mnc,
		Tpt:  s.transport,
	}, nil
}

func (s *UtpSocket) Matches(a ma.Multiaddr) bool {
	return mafmt.UTP.Matches(a)
}

func (t *UtpSocket) Close() error {
	return t.s.Close()
}

func (t *UtpSocket) Addr() net.Addr {
	return t.s.Addr()
}

func (t *UtpSocket) Multiaddr() ma.Multiaddr {
	return t.laddr
}

var _ tpt.Transport = (*UtpTransport)(nil)

func parseUtpNetAddr(a net.Addr) (ma.Multiaddr, error) {
	var udpaddr *net.UDPAddr
	switch a := a.(type) {
	case *utp.Addr:
		udpaddr = a.Child.(*net.UDPAddr)
	case *net.UDPAddr:
		udpaddr = a
	default:
		return nil, fmt.Errorf("was not given a valid utp address")
	}

	// Get IP Addr
	ipm, err := manet.FromIP(udpaddr.IP)
	if err != nil {
		return nil, errIncorrectNetAddr
	}

	// Get UDP Addr
	utpm, err := ma.NewMultiaddr(fmt.Sprintf("/udp/%d/utp", udpaddr.Port))
	if err != nil {
		return nil, errIncorrectNetAddr
	}

	// Encapsulate
	return ipm.Encapsulate(utpm), nil
}

func parseUtpMaddr(maddr ma.Multiaddr) (net.Addr, error) {
	utpbase, err := ma.NewMultiaddr("/utp")
	if err != nil {
		return nil, err
	}

	raw := maddr.Decapsulate(utpbase)

	udpa, err := manet.ToNetAddr(raw)
	if err != nil {
		return nil, err
	}

	return &utp.Addr{udpa}, nil
}
