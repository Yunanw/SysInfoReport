package SysInfo

import "github.com/shirou/gopsutil/v3/mem"

type MemoryInfo struct {
	Total       float64
	Used        float64
	UsedPercent float64
}

func collectMemory(sysInfo *SysInfo) {
	v, _ := mem.VirtualMemory()

	sysInfo.Memory.Total = Round(ToGB(v.Total))
	sysInfo.Memory.Used = Round(ToMB(v.Used))
	sysInfo.Memory.UsedPercent = v.UsedPercent
}
