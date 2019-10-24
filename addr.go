package utp

// This file manage the parsing and registering part.

import (
	"fmt"
	"net"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

var utpAddrSpec = &manet.NetCodec{
	ProtocolName:     "utp",
	NetAddrNetworks:  []string{"utp", "utp4", "utp6", "utp/udp", "utp/udp4", "utp/udp6"},
	ParseNetAddr:     parseUtpNetAddr,
	ConvertMultiaddr: parseUtpMaddr,
}

func init() {
	manet.RegisterNetCodec(utpAddrSpec)
}

func parseUtpNetAddr(udpaddr net.Addr) (ma.Multiaddr, error) {
	switch udpaddr := udpaddr.(type) {
	case *net.UDPAddr:
		// Get IP Addr
		ipm, err := manet.FromIP(udpaddr.IP)
		if err != nil {
			return nil, err
		}

		// Get UDP Addr
		utpm, err := ma.NewMultiaddr(fmt.Sprintf("/udp/%d/utp", udpaddr.Port))
		if err != nil {
			return nil, err
		}

		// Encapsulate
		return ipm.Encapsulate(utpm), nil
	default:
		return nil, fmt.Errorf("\"%#v\" is not given a valid utp address", udpaddr)
	}

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

	return udpa, nil
}
