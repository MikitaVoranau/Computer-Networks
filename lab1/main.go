package main

import (
	"fmt"
	"net"
)

func main() {

	netInterfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, iFace := range netInterfaces {

		if !isInterfaceCorrect(iFace) {
			continue
		}

		fmt.Println("Имя сетевого интерфейса - ", iFace.Name, "\nЕго MAC адрес - ", iFace.HardwareAddr)

		adds, err := iFace.Addrs()
		if err != nil {
			_ = fmt.Errorf("не удалось получить адреса %w для сетевого интерфейса %s", err, iFace.Name)
		}

		getIPAdds(adds)
	}
}

func getIPAdds(adds []net.Addr) {
	for _, addr := range adds {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipNet.IP.To4() != nil {
			fmt.Printf("IPv4-адрес: %s\n", ipNet)

			startIP, endIP := CalculatesIPRange(ipNet)
			fmt.Printf("Диапазон IPv4-адресов: %s - %s\n", startIP, endIP)
			//SendPings(startIP)
		}
	}
}
