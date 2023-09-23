package lib

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	sigar "github.com/cloudfoundry/gosigar"
)

// SystemInfo contains basic information about system load
type SystemInfo struct {
	Memory      sigar.Mem
	Swap        sigar.Swap
	Uptime      int
	UptimeS     string
	LoadAvg     sigar.LoadAverage
	CPUList     sigar.CpuList
	Arch        string
	Os          string
	CurrentTime time.Time
}

// GetSystemInfo returns short info about system load
func GetSystemInfo() SystemInfo {
	s := SystemInfo{}

	uptime := sigar.Uptime{}
	if err := uptime.Get(); err == nil {
		s.Uptime = int(uptime.Length)
		s.UptimeS = uptime.Format()
	}

	avg := sigar.LoadAverage{}
	if err := avg.Get(); err == nil {
		s.LoadAvg = avg
		s.LoadAvg.One = formatFloat(s.LoadAvg.One)
		s.LoadAvg.Five = formatFloat(s.LoadAvg.Five)
		s.LoadAvg.Fifteen = formatFloat(s.LoadAvg.Fifteen)
	}

	s.CurrentTime = time.Now()

	mem := sigar.Mem{}
	if err := mem.Get(); err == nil {
		s.Memory = mem
	}

	swap := sigar.Swap{}
	if err := swap.Get(); err == nil {
		s.Swap = swap
	}

	cpulist := sigar.CpuList{}
	if err := cpulist.Get(); err == nil {
		s.CPUList = cpulist
	}

	s.Arch = runtime.GOARCH
	s.Os = runtime.GOOS

	return s
}

func formatFloat(f float64) float64 {
	formatted, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
	return formatted
}
