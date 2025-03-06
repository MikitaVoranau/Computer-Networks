package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	hostIP, _ := getHostInfo()

	ourIface := getIface()
	adds, err := ourIface.Addrs()
	if err != nil {
		log.Fatalf("не удалось получить адреса %w для сетевого интерфейса %s", err, ourIface.Name)
	}
	fmt.Printf("Собственный компьютер - IP: %s, MAC: %s\n\n", hostIP, ourIface.HardwareAddr.String())
	Ip4Adds := getIPAdds(adds)
	if Ip4Adds == nil {
		log.Fatalf("couldn't get an ip4 address from  your iface ")
	}

	startIP, endIP := CalculatesIPRange(Ip4Adds)
	results := make(chan DeviceInfo)

	go func() {
		for ip := range results {
			fmt.Println("Доступное устройство:", ip)
		}
	}()
	SendPingstoIPs(startIP, endIP, results, hostIP)

	close(results)

}

func getIPAdds(adds []net.Addr) *net.IPNet {
	for _, address := range adds {
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			log.Fatalf("getIPAdds : cannot type to net")
		}
		if ipNet.IP.To4() != nil {
			return ipNet
		}
	}
	return nil
}
