package collector

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/mem"
)

type MemoryInfo struct {
	Total       float64
	Used        float64
	UsedPercent float64
}

type MemoryCollector struct {
}

func (collector *MemoryCollector) Collect(sysInfo *SysInfo) error {
	v, err := mem.VirtualMemory()
	if err != nil {
		return errors.Wrap(err, "读取系统内存信息失败")
	}
	sysInfo.Memory.Total = Round(ToGB(v.Total))
	sysInfo.Memory.Used = Round(ToMB(v.Used))
	sysInfo.Memory.UsedPercent = v.UsedPercent

	return nil
}
