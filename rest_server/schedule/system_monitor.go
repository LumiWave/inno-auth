package schedule

import (
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/config"
	"github.com/ONBUFF-IP-TOKEN/inno-auth/rest_server/controllers/context"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

var TagVersion = "v1.0.1"

var gSystemMonitor *SystemMonitor
var onceSystemMonitor sync.Once

type SystemMonitor struct {
	NodeMetric *context.NodeMetric
}

func GetSystemMonitor() *SystemMonitor {
	onceSystemMonitor.Do(func() {
		gSystemMonitor = new(SystemMonitor)
		gSystemMonitor.NodeMetric = new(context.NodeMetric)
		gSystemMonitor.NodeMetric.UpTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
		gSystemMonitor.Run()
	})

	return gSystemMonitor
}

func (o *SystemMonitor) Run() {
	go func() {
		ticker := time.NewTicker(time.Duration(10) * time.Second)

		for {
			o.CheckMetricInfo()
			<-ticker.C
		}
	}()
}

func (o *SystemMonitor) GetMetricInfo() *context.NodeMetric {
	return o.NodeMetric
}

func (o *SystemMonitor) CheckMetricInfo() *context.NodeMetric {
	conf := config.GetInstance()

	// host
	o.NodeMetric.Host = conf.InnoAuth.ApplicationName
	// 버전
	o.NodeMetric.Version = TagVersion
	// 정상여부
	o.NodeMetric.IsRunning = true
	// cpu 타임
	o.NodeMetric.CpuTime = strconv.FormatInt(time.Now().UnixMilli(), 10)
	// 메모리 사용정보
	o.NodeMetric.MemTotalBytes, o.NodeMetric.MemAllocBytes, o.NodeMetric.MemPercent = o.getMemoryUsage()
	// CPU 점유율
	o.NodeMetric.CpuUsage = int32(o.getCpuUsage())
	// disk 사용 정보
	o.NodeMetric.DiskUsage = o.getDisUsage()

	return o.NodeMetric
}

func (o *SystemMonitor) getCpuUsage() uint64 {
	percent, err := cpu.Percent(0, false)
	if err != nil || len(percent) <= 0 {
		return 0
	}

	return uint64(math.Ceil(percent[0]))
}

func (o *SystemMonitor) getMemoryUsage() (uint64, uint64, float32) {
	vmStat, err := mem.VirtualMemory()
	if err != nil || vmStat == nil {
		return 0, 0, 0
	}

	return vmStat.Total, vmStat.Used, float32(vmStat.Used) * 100 / float32(vmStat.Total)
}

func (o *SystemMonitor) getDisUsage() []context.DiskUsage {
	disks := []context.DiskUsage{}
	if partitions, err := disk.Partitions(true); err != nil {
		return disks
	} else {
		for _, partition := range partitions {
			// 주요 디스크만 수집
			if strings.Index(partition.Mountpoint, "/run") == 0 ||
				strings.Index(partition.Mountpoint, "/boot") == 0 ||
				strings.EqualFold(partition.Mountpoint, "/") ||
				strings.Index(partition.Mountpoint, "/mnt/") == 0 ||
				strings.Index(partition.Mountpoint, "C") == 0 ||
				strings.Index(partition.Mountpoint, "D") == 0 {
				if state, err := disk.Usage(partition.Mountpoint); err == nil {
					disks = append(disks, context.DiskUsage{Disk: *state})
				}
			}

		}
	}

	return disks
}
