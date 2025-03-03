package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"net"
)

func SendPings(allIps []net.IP) []net.IP {
	for _, value := range allIps {
		pinger, err := ping.NewPinger(net.IP.String(value))
		if err != nil {
			panic(err)
		}
		pinger.Count = 3
		err = pinger.Run() // Blocks until finished.
		if err != nil {
			panic(err)
		}
		stats := pinger.Statistics()
		fmt.Println(stats)
	}
	return allIps
}
