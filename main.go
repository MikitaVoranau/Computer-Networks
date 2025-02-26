package main

import (
	"fmt"
	"net"
)

func main() {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		_ = fmt.Errorf("ошибка при получении интерфейсов - %w", err)
		return
	}
	for _, iface := range netInterfaces {
		if (iface.Flags&net.FlagUp == 0) || iface.HardwareAddr == nil {
			continue
		}
		fmt.Println("Имя сетевого интерфейса - ", iface.Name, "\nЕго MAC адрес - ", iface.HardwareAddr)
		addrs, err := iface.Addrs()
		if err != nil {
			_ = fmt.Errorf("не удалось получить адреса %w для сетевого интерфейса %s", err, iface.Name)
		}
		if len(addrs) > 0 {
			for _, addr := range addrs {
				fmt.Printf("IP-адрес: %s\n", addr)
			}
			fmt.Println("-----")
		}

	}
}
