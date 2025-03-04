package main

import (
	"fmt"
	"log"
	"net"
)

func main() {

	ourIface := getIface()

	adds, err := ourIface.Addrs()
	if err != nil {
		log.Fatalf("не удалось получить адреса %w для сетевого интерфейса %s", err, ourIface.Name)
	}

	Ip4Adds := getIPAdds(adds)
	if Ip4Adds == nil {
		log.Fatalf("couldn't get an ip4 address from  your iface ")
	}

	startIP, endIP := CalculatesIPRange(Ip4Adds)
	results := make(chan string)

	go func() {
		for ip := range results {
			fmt.Println("Доступное устройство:", ip)
		}
	}()

	// Запускаем сканирование
	SendPingstoIPs(startIP, endIP, results)

	// Закрываем канал после завершения сканирования
	close(results)

}

func getIPAdds(adds []net.Addr) *net.IPNet {
	for _, address := range adds {
		ipNet, ok := address.(*net.IPNet)
		if !ok {
			log.Fatalf("getIPAdds : cannot type to net")
		}
		if ipNet.IP.To4() != nil {
			fmt.Printf("IPv4-адрес: %s\n", ipNet)
			return ipNet
		}
	}
	return nil
}
