package main

import (
	"github.com/go-ping/ping"
	"log"
	"net"
	"sync"
	"time"
)

type Semaphore struct {
	semaCh chan struct{}
}

var (
	wg      sync.WaitGroup
	mutex   sync.Mutex
	results = make(chan string)
)

const maxGoroutines = 200

func NewSemaphore(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.semaCh
}

func SendPingstoIPs(startIP, endIP net.IP, results chan<- string) {
	sem := NewSemaphore(200)

	for ip := startIP; !ip.Equal(endIP); IncourIP(ip) {

		wg.Add(1)
		sem.Acquire() // Забираем слот

		go func(ip string) {
			defer wg.Done()
			defer sem.Release() // Освобождаем слот
			pingIP(ip, results)
		}(ip.String())
	}

	wg.Wait()
}

func pingIP(ip string, results chan<- string) {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		log.Printf("Ошибка при создании пингера для %s: %v\n", ip, err)
		return
	}

	pinger.SetPrivileged(true)
	pinger.Count = 1
	pinger.Timeout = time.Second * 5
	pinger.Size = 600

	err = pinger.Run()
	if err != nil {
		log.Printf("Ошибка при пинге %s: %v\n", ip, err)
		return
	}

	stats := pinger.Statistics()
	if stats.PacketsRecv > 0 {
		mutex.Lock()
		results <- ip
		mutex.Unlock()
	}
}

func IncourIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
