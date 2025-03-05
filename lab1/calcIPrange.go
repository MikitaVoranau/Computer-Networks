package main

import "net"

func CalculatesIPRange(ipNet *net.IPNet) (net.IP, net.IP) {

	startIP := ipNet.IP.To4()
	mask := ipNet.Mask

	networkIp := make(net.IP, len(startIP))
	for i := range startIP {
		networkIp[i] = startIP[i] & mask[i]

	}

	endIP := make(net.IP, len(startIP))

	for i := range startIP {
		endIP[i] = networkIp[i] | ^mask[i]

	}
	return IncIP(networkIp, false), endIP
}

func GetIPRange(startIP, endIP *net.IP) []net.IP {
	var ips []net.IP

	for ip := *startIP; !ip.Equal(*endIP); ip = IncIP(ip, false) {
		ips = append(ips, ip)
	}
	ips = append(ips, *endIP)
	return ips
}
