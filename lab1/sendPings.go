package main

import (
	"github.com/go-ping/ping"
	"log"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type Semaphore struct {
	semaCh chan struct{}
}

type DeviceInfo struct {
	IP   string
	MAC  string
	Name string
}

var (
	wg      sync.WaitGroup
	mutex   sync.Mutex
	results = make(chan DeviceInfo) // Канал для передачи информации об устройстве
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

func SendPingstoIPs(startIP, endIP net.IP, results chan<- DeviceInfo) {
	sem := NewSemaphore(maxGoroutines)

	for ip := startIP; !ip.Equal(endIP); IncIP(ip, true) {
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

func pingIP(ip string, results chan<- DeviceInfo) {
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
		mac, name := getARPInfo(ip) // Получаем MAC-адрес и имя устройства
		results <- DeviceInfo{IP: ip, MAC: mac, Name: name}
	}
}

func getARPInfo(ip string) (string, string) {
	cmd := exec.Command("arp", "-a")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Ошибка при выполнении команды arp -a: %v\n", err)
		return "", ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				mac := parts[1]
				name := ""
				if len(parts) >= 4 {
					name = parts[3]
				}
				return mac, name
			}
		}
	}

	return "", ""
}

func IncIP(ip net.IP, inPlace bool) net.IP {
	var targetIP net.IP
	if inPlace {
		targetIP = ip
	} else {
		targetIP = make(net.IP, len(ip))
		copy(targetIP, ip)
	}

	for i := len(targetIP) - 1; i >= 0; i-- {
		targetIP[i]++
		if targetIP[i] != 0 {
			break
		}
	}

	if !inPlace {
		return targetIP
	}
	return nil
}
