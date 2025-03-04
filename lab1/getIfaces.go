package main

import (
	"fmt"
	"log"
	"net"
)

func getIface() net.Interface {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("getIface : cannot get interfaces %w", err)
	}

	netInterfaces = findActiveIfaces(netInterfaces)

	for i, Iface := range netInterfaces {
		fmt.Println(i+1, ")", Iface.Name)
	}

	var choose int
	_, err = fmt.Scan(&choose)
	if err != nil {
		log.Fatalf("getIface : wrong choose %w", err)
	}
	if choose < 0 || choose > len(netInterfaces) {
		log.Fatalf("no such interface")
	}
	return netInterfaces[choose-1]
}

func findActiveIfaces(allIfaces []net.Interface) []net.Interface {

	var activeIFaces []net.Interface

	for _, iface := range allIfaces {

		if iface.Flags&net.FlagUp == 0 || iface.HardwareAddr == nil {
			continue

		}
		activeIFaces = append(activeIFaces, iface)
	}

	return activeIFaces
}
