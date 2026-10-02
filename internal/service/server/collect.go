package server

import (
	"math"
	"net"
	"os"
	"runtime"
	"syscall"
	"time"
)

// platformMetrics 由各平台实现填充（见 metrics_linux.go / metrics_darwin.go）
type platformMetrics struct {
	CPUPercent float64
	CPUModel   string
	Kernel     string
	Platform   string
	Uptime     uint64
	LoadAvg    float64
	MemTotal   uint64
	MemUsed    uint64
	SwapTotal  uint64
	SwapUsed   uint64
}

// collect 汇总主机信息。指标按“尽力而为”采集，单项取不到时留零值，不阻断整个接口。
// 目前仅实现 linux / darwin：两者都能通过 syscall.Statfs 取磁盘，其余指标见各自构建文件。
func collect() *Info {
	pm := collectPlatform()

	info := &Info{
		Hostname:    hostname(),
		IP:          localIP(),
		OS:          runtime.GOOS,
		Platform:    pm.Platform,
		Kernel:      pm.Kernel,
		Arch:        runtime.GOARCH,
		CPUModel:    pm.CPUModel,
		CPUCores:    runtime.NumCPU(),
		Uptime:      pm.Uptime,
		LoadAvg:     round1(pm.LoadAvg),
		CPUPercent:  clampPercent(pm.CPUPercent),
		MemTotal:    pm.MemTotal,
		MemUsed:     pm.MemUsed,
		MemPercent:  percent(pm.MemUsed, pm.MemTotal),
		SwapTotal:   pm.SwapTotal,
		SwapUsed:    pm.SwapUsed,
		SwapPercent: percent(pm.SwapUsed, pm.SwapTotal),
		CollectedAt: time.Now().Unix(),
	}

	info.DiskTotal, info.DiskUsed = diskUsage("/")
	info.DiskPercent = percent(info.DiskUsed, info.DiskTotal)
	return info
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return ""
	}
	return name
}

// localIP 优先返回内网 IPv4，取不到时退回第一个非回环 IPv4
func localIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	fallback := ""
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}
			if ip.IsPrivate() {
				return ip.String()
			}
			if fallback == "" {
				fallback = ip.String()
			}
		}
	}
	return fallback
}

// diskUsage 取指定挂载点的总容量与已用容量
func diskUsage(path string) (total, used uint64) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(path, &stat); err != nil {
		return 0, 0
	}
	blockSize := uint64(stat.Bsize)
	total = stat.Blocks * blockSize
	available := stat.Bavail * blockSize
	return total, total - available
}

func percent(used, total uint64) float64 {
	if total == 0 {
		return 0
	}
	return round1(float64(used) / float64(total) * 100)
}

func round1(value float64) float64 {
	return math.Round(value*10) / 10
}

func clampPercent(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 100 {
		return 100
	}
	return round1(value)
}
