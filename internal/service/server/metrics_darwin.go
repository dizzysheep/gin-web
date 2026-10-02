//go:build darwin

package server

import (
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// darwinScript 一次 sh 调用取回全部指标：macOS 没有 /proc，只能依赖系统命令。
// top 需要两次采样（间隔 1s）才能得到瞬时 CPU 占用，所以这里耗时约 1s，
// 服务层的 3s 缓存保证前端轮询不会反复 fork。
const darwinScript = `
echo "### vm_stat"
vm_stat
echo "### mem_total"
sysctl -n hw.memsize
echo "### swap"
sysctl -n vm.swapusage
echo "### cpu_model"
sysctl -n machdep.cpu.brand_string
echo "### hw_model"
sysctl -n hw.model
echo "### cpu_usage"
top -l 2 -n 0 -s 1 | grep 'CPU usage'
echo "### boottime"
sysctl -n kern.boottime
echo "### loadavg"
sysctl -n vm.loadavg
echo "### kernel"
uname -r
echo "### platform"
sw_vers -productName
sw_vers -productVersion
`

func collectPlatform() platformMetrics {
	sections := readSections()

	memTotal := parseUint(firstLine(sections["mem_total"]))
	memUsed := parseVMStatUsed(sections["vm_stat"])
	swapLine := firstLine(sections["swap"])
	swapTotal := parseSwapValue(swapLine, "total")
	swapUsed := parseSwapValue(swapLine, "used")
	loadAvg := parseLoadAvg(firstLine(sections["loadavg"]))

	metrics := platformMetrics{
		CPUModel:  firstNonEmpty(firstLine(sections["cpu_model"]), firstLine(sections["hw_model"])),
		Kernel:    firstLine(sections["kernel"]),
		Platform:  strings.TrimSpace(strings.Join(sections["platform"], " ")),
		Uptime:    parseUptime(sections["boottime"]),
		LoadAvg:   loadAvg,
		MemTotal:  memTotal,
		MemUsed:   memUsed,
		SwapTotal: swapTotal,
		SwapUsed:  swapUsed,
	}

	metrics.CPUPercent = parseCPUUsage(sections["cpu_usage"])
	if metrics.CPUPercent < 0 {
		// top 取不到时用 1 分钟负载均值折算，至少反映趋势
		metrics.CPUPercent = loadAvg * 100 / float64(runtime.NumCPU())
	}
	return metrics
}

// readSections 执行采集脚本并按 "### 名称" 切成小节
func readSections() map[string][]string {
	output, err := exec.Command("/bin/sh", "-c", darwinScript).Output()
	if err != nil && len(output) == 0 {
		return nil
	}

	sections := make(map[string][]string)
	current := ""
	for _, line := range strings.Split(string(output), "\n") {
		if name, ok := strings.CutPrefix(line, "### "); ok {
			current = strings.TrimSpace(name)
			continue
		}
		if current == "" {
			continue
		}
		if line = strings.TrimSpace(line); line != "" {
			sections[current] = append(sections[current], line)
		}
	}
	return sections
}

// parseVMStatUsed 用 active + wired + 压缩内存估算已用内存，与活动监视器口径接近
func parseVMStatUsed(lines []string) uint64 {
	var pageSize, active, wired, compressed uint64
	for _, line := range lines {
		if index := strings.Index(line, "page size of"); index >= 0 {
			pageSize = parseUint(firstField(strings.TrimPrefix(line[index:], "page size of")))
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		pages := parseUint(strings.Trim(value, ". "))
		switch strings.TrimSpace(key) {
		case "Pages active":
			active = pages
		case "Pages wired down":
			wired = pages
		case "Pages occupied by compressor":
			compressed = pages
		}
	}
	return (active + wired + compressed) * pageSize
}

// parseSwapValue 解析 "total = 4096.00M  used = 1234.50M  free = ..."
func parseSwapValue(line, key string) uint64 {
	_, rest, ok := strings.Cut(line, key+" =")
	if !ok {
		return 0
	}
	return parseSize(firstField(rest))
}

// parseSize 解析 "2048.00M" / "512K" / "2.00G" 形式的容量
func parseSize(value string) uint64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	unit := uint64(1)
	switch value[len(value)-1] {
	case 'K', 'k':
		unit = 1024
	case 'M', 'm':
		unit = 1024 * 1024
	case 'G', 'g':
		unit = 1024 * 1024 * 1024
	default:
		size, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0
		}
		return uint64(size)
	}

	size, err := strconv.ParseFloat(value[:len(value)-1], 64)
	if err != nil {
		return 0
	}
	return uint64(size * float64(unit))
}

// parseCPUUsage 取最后一次采样，解析 "CPU usage: 5.55% user, 10.0% sys, 84.44% idle"；
// 解析失败返回 -1 交由调用方回退
func parseCPUUsage(lines []string) float64 {
	for index := len(lines) - 1; index >= 0; index-- {
		line := lines[index]
		position := strings.Index(line, "% idle")
		if position < 0 {
			continue
		}
		idleText := line[:position]
		if comma := strings.LastIndex(idleText, ","); comma >= 0 {
			idleText = idleText[comma+1:]
		}
		idle, err := strconv.ParseFloat(strings.TrimSpace(idleText), 64)
		if err != nil {
			continue
		}
		return 100 - idle
	}
	return -1
}

// parseUptime 由 kern.boottime 的 "{ sec = 1750000000, usec = ... }" 推算运行时长
func parseUptime(lines []string) uint64 {
	for _, line := range lines {
		_, rest, ok := strings.Cut(line, "sec =")
		if !ok {
			continue
		}
		value, _, _ := strings.Cut(rest, ",")
		bootAt := parseUint(strings.TrimSpace(value))
		if bootAt == 0 {
			continue
		}
		now := uint64(time.Now().Unix())
		if now > bootAt {
			return now - bootAt
		}
	}
	return 0
}

// parseLoadAvg 解析 "{ 1.23 4.56 7.89 }" 中的 1 分钟负载
func parseLoadAvg(line string) float64 {
	fields := strings.Fields(strings.Trim(line, "{} "))
	if len(fields) == 0 {
		return 0
	}
	value, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0
	}
	return value
}

func parseUint(value string) uint64 {
	parsed, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0
	}
	return parsed
}

func firstField(value string) string {
	fields := strings.Fields(value)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func firstLine(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return lines[0]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
