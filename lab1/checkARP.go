package main

import (
	"log"
	"os/exec"
	"strings"
)

func GetARPInfo(ip string) string {
	cmd := exec.Command("arp", "-a", ip)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Ошибка при выполнении команды arp -a: %v\n", err)
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, ip) {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				mac := parts[1]
				return mac
			}
		}
	}

	return ""
}
