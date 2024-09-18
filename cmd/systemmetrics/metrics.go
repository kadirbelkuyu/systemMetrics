package main

import (
	"fmt"
	"log"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func main() {
	// System
	hostInfo, err := host.Info()
	if err != nil {
		log.Fatalf("Host bilgileri alınırken hata: %v", err)
	}
	fmt.Printf("Host Bilgileri: %s %s (Uptime: %d saniye)\n\n", hostInfo.Platform, hostInfo.PlatformVersion, hostInfo.Uptime)

	// CPU
	cpuInfo, err := cpu.Info()
	if err != nil {
		log.Fatalf("CPU bilgileri alınırken hata: %v", err)
	}
	for _, ci := range cpuInfo {
		fmt.Printf("CPU: %s %s\n", ci.ModelName, ci.Cores)
	}

	// CPU data
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		log.Fatalf("CPU kullanımı alınırken hata: %v", err)
	}
	fmt.Printf("CPU Kullanımı: %.2f%%\n", cpuPercent[0])

	// Memory
	virtualMem, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Bellek bilgileri alınırken hata: %v", err)
	}
	fmt.Printf("Toplam Bellek: %v MB\nKullanılan Bellek: %v MB (%.2f%%)\n\n", virtualMem.Total/1024/1024, virtualMem.Used/1024/1024, virtualMem.UsedPercent)

	// Disk
	diskInfo, err := disk.Usage("/")
	if err != nil {
		log.Fatalf("Disk bilgileri alınırken hata: %v", err)
	}
	fmt.Printf("Disk Toplam: %v GB\nDisk Kullanılan: %v GB (%.2f%%)\n", diskInfo.Total/1024/1024/1024, diskInfo.Used/1024/1024/1024, diskInfo.UsedPercent)
	
	for {
		cpuPercent, _ := cpu.Percent(1*time.Second, false)
		virtualMem, _ := mem.VirtualMemory()

		fmt.Printf("\nAnlık CPU Kullanımı: %.2f%%\n", cpuPercent[0])
		fmt.Printf("Anlık Bellek Kullanımı: %v MB (%.2f%%)\n", virtualMem.Used/1024/1024, virtualMem.UsedPercent)
		time.Sleep(5 * time.Second)
	}
}
