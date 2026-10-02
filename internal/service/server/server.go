package server

import (
	"sync"
	"time"
)

// cacheTTL 指标缓存时长。darwin 上采集 CPU 需要执行 top（约 1s），
// 缓存后前端轮询不会反复 fork 进程，CPU 差值也有了稳定的采样窗口。
const cacheTTL = 3 * time.Second

// Info 服务器信息与资源占用（占用率单位为百分比，容量单位为字节）
type Info struct {
	Hostname string  `json:"hostname"`
	IP       string  `json:"ip"`
	OS       string  `json:"os"`
	Platform string  `json:"platform"`
	Kernel   string  `json:"kernel"`
	Arch     string  `json:"arch"`
	CPUModel string  `json:"cpu_model"`
	CPUCores int     `json:"cpu_cores"`
	Uptime   uint64  `json:"uptime"`
	LoadAvg  float64 `json:"load_avg"`

	CPUPercent float64 `json:"cpu_percent"`

	MemTotal   uint64  `json:"mem_total"`
	MemUsed    uint64  `json:"mem_used"`
	MemPercent float64 `json:"mem_percent"`

	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
	SwapPercent float64 `json:"swap_percent"`

	DiskTotal   uint64  `json:"disk_total"`
	DiskUsed    uint64  `json:"disk_used"`
	DiskPercent float64 `json:"disk_percent"`

	CollectedAt int64 `json:"collected_at"`
}

type ServerService interface {
	Info() *Info
}

type serverService struct {
	mu      sync.Mutex
	cached  *Info
	cacheAt time.Time
}

func NewServerService() ServerService {
	return &serverService{}
}

// Info 返回后端进程所在主机的实时信息，结果在进程内缓存 cacheTTL
func (s *serverService) Info() *Info {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cached != nil && time.Since(s.cacheAt) < cacheTTL {
		return s.cached
	}

	s.cached = collect()
	s.cacheAt = time.Now()
	return s.cached
}
