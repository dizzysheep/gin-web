//go:build linux

package server

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	cpuMu     sync.Mutex
	cpuPrevAt time.Time
	cpuIdle   uint64
	cpuTotal  uint64
)

// collectPlatform 通过 /proc 读取 Linux 主机指标
func collectPlatform() platformMetrics {
	memTotal, memAvailable, swapTotal, swapFree := readMemInfo()

	return platformMetrics{
		CPUPercent: cpuPercent(),
		CPUModel:   readCPUModel(),
		Kernel:     readTrimmed("/proc/sys/kernel/osrelease"),
		Platform:   readOSRelease(),
		Uptime:     readUptime(),
		LoadAvg:    readLoadAvg(),
		MemTotal:   memTotal,
		MemUsed:    memTotal - memAvailable,
		SwapTotal:  swapTotal,
		SwapUsed:   swapTotal - swapFree,
	}
}

// cpuPercent 由两次 /proc/stat 采样求差值；距上次采样过近时补采一次
func cpuPercent() float64 {
	idle, total := readCPUTimes()
	now := time.Now()

	cpuMu.Lock()
	defer cpuMu.Unlock()

	if cpuTotal == 0 || now.Sub(cpuPrevAt) < 200*time.Millisecond {
		time.Sleep(200 * time.Millisecond)
		idle, total = readCPUTimes()
		now = time.Now()
	}

	var value float64
	if total > cpuTotal {
		deltaTotal := total - cpuTotal
		deltaIdle := idle - cpuIdle
		value = float64(deltaTotal-deltaIdle) / float64(deltaTotal) * 100
	}

	cpuIdle, cpuTotal, cpuPrevAt = idle, total, now
	return value
}

// readCPUTimes 聚合 /proc/stat 首个 cpu 行，idle 含 idle 与 iowait 两列
func readCPUTimes() (idle, total uint64) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0
	}

	fields := []string(nil)
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "cpu ") {
			fields = strings.Fields(line)
			break
		}
	}
	if len(fields) < 5 {
		return 0, 0
	}

	for index, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			continue
		}
		if index == 3 || index == 4 {
			idle += value
		}
		total += value
	}
	return idle, total
}

func readMemInfo() (total, available, swapTotal, swapFree uint64) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return 0, 0, 0, 0
	}

	var free uint64
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			continue
		}
		value *= 1024 // /proc/meminfo 单位为 kB

		switch strings.TrimSuffix(fields[0], ":") {
		case "MemTotal":
			total = value
		case "MemAvailable":
			available = value
		case "MemFree":
			free = value
		case "SwapTotal":
			swapTotal = value
		case "SwapFree":
			swapFree = value
		}
	}
	if available == 0 {
		available = free
	}
	return total, available, swapTotal, swapFree
}

func readUptime() uint64 {
	fields := strings.Fields(readTrimmed("/proc/uptime"))
	if len(fields) == 0 {
		return 0
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil || seconds < 0 {
		return 0
	}
	return uint64(seconds)
}

func readLoadAvg() float64 {
	fields := strings.Fields(readTrimmed("/proc/loadavg"))
	if len(fields) == 0 {
		return 0
	}
	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return value
}

// readOSRelease 取 /etc/os-release 的 PRETTY_NAME，如 "Ubuntu 22.04.4 LTS"
func readOSRelease() string {
	data, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		value, ok := strings.CutPrefix(line, "PRETTY_NAME=")
		if !ok {
			continue
		}
		return strings.Trim(strings.TrimSpace(value), `"`)
	}
	return ""
}

// readCPUModel x86 取 model name，arm64 取 Hardware / Model
func readCPUModel() string {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return ""
	}

	fallback := ""
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		value = strings.TrimSpace(value)
		switch strings.TrimSpace(key) {
		case "model name":
			return value
		case "Hardware", "Model":
			if fallback == "" {
				fallback = value
			}
		}
	}
	return fallback
}

func readTrimmed(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
