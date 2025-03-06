package main

import (
	"log"
	"net"
)

func getHostInfo() (string, string) {
	hostIP := getOutboundIP()

	mac := GetARPInfo(hostIP)
	return hostIP, mac
}

func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Ошибка при получении IP-адреса: %v", err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
