package SysInfo

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/thoas/go-funk"
)

type CPUInfo struct {
	CpuCount    int
	LoadPercent []float64
}

func collectCPU(sysInfo *SysInfo) {
	sysInfo.CPU.LoadPercent, _ = cpu.Percent(0, true)
	sysInfo.CPU.LoadPercent = funk.Map(sysInfo.CPU.LoadPercent, sysInfo.round).([]float64)
	sysInfo.CPU.CpuCount = len(sysInfo.CPU.LoadPercent)
}
