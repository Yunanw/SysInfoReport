package SysInfo

import (
	. "github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
)

type CPUInfo struct {
	CpuCount    int
	LoadPercent []float64
}

func collectCPU(sysInfo *SysInfo) error {
	p, err := cpu.Percent(0, true)
	if err != nil {
		return errors.Wrap(err, "读取CPU使用信息失败")
	}

	From(p).SelectT(func(f float64) float64 {
		return Round(f)
	}).ToSlice(&sysInfo.CPU.LoadPercent)
	sysInfo.CPU.CpuCount = len(sysInfo.CPU.LoadPercent)
	return nil
}
