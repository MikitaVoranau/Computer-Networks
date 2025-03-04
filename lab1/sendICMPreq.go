package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"net"
	"sync"
	"time"
)

var wg sync.WaitGroup

var results = make(chan string)

func SendPingstoIPs(startIP, endIP net.IP) {
	go func() {
		for result := range results {
			fmt.Println(result)
		}
	}()

	for ip := startIP; !ip.Equal(endIP); IncIP(ip) {
		wg.Add(1)
		pingIP(ip.String(), &wg, results)
	}

	wg.Wait()

	close(results)
}

func pingIP(ip string, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Fatalf("Ping IP : cannot send ping %w", err)
	}
	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = time.Second * 1
	err = pinger.Run()
	if err != nil {
		log.Fatalf("Ping IP : cannot finished ping, %w", err)
	}
	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		results <- ip
	}
	fmt.Println(stats)
}
