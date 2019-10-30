package utp

// This file manage the parsing and registering part.

import (
	"net"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
)

var utpAddrSpec = &manet.NetCodec{
	ProtocolName: "utp",
	NetAddrNetworks: []string{"utp", "utp4", "utp6",
		// Theses are deaprecated and only here for compatibilty issue, use the top ones.
		"utp/udp", "utp/udp4", "utp/udp6"},
	ParseNetAddr:     parseUtpNetAddr,
	ConvertMultiaddr: parseUtpMaddr,
}

func init() {
	manet.RegisterNetCodec(utpAddrSpec)
}

func parseUtpNetAddr(addr net.Addr) (ma.Multiaddr, error) {
	addrStr := addr.String()
	portSep := strings.LastIndex(addrStr, ":")
	ip := addrStr[0:portSep]
	port := addrStr[portSep+1 : len(addrStr)]
	var s string
	if strings.Contains(addr.Network(), "6") || strings.Contains(ip, "[") || strings.Contains(ip, "]") {
		s = "/ip6/" + ip + "/udp/" + port
	} else {
		s = "/ip4/" + ip + "/udp/" + port
	}

	mabase, err := ma.NewMultiaddr(s)
	if err != nil {
		return nil, err
	}
	return mabase.Encapsulate(emptyUtpMa), nil
}

func parseUtpMaddr(maddr ma.Multiaddr) (net.Addr, error) {
	return manet.ToNetAddr(maddr.Decapsulate(emptyUtpMa))
}
