package context

import "github.com/shirou/gopsutil/disk"

type DiskUsage struct {
	Disk disk.UsageStat
}

type NodeMetric struct {
	Host string `json:"host"`

	Version       string      `json:"version"`
	IsRunning     bool        `json:"is_running"`
	UpTime        string      `json:"up_time"`
	CpuTime       string      `json:"cpu_time"`
	MemTotalBytes uint64      `json:"mem_total_bytes"`
	MemAllocBytes uint64      `json:"mem_alloc_bytes"`
	MemPercent    float32     `json:"mem_usage_percent"`
	CpuUsage      int32       `json:"cpu_usage"`
	DiskUsage     []DiskUsage `json:"disk_usage"`
}
