package main

import (
	"log"
	"net"
)

func getHostInfo() (string, string) {
	// Получаем IP-адрес
	hostIP := getOutboundIP()

	// Получаем MAC-адрес и имя
	mac := GetARPInfo(hostIP)
	return hostIP, mac
}

// getOutboundIP возвращает IP-адрес собственного компьютера
func getOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("Ошибка при получении IP-адреса: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}
