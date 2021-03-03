package SysInfo

import (
	"SysInfoReport/pkg/config"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/thoas/go-funk"
)

type ProcessInfo struct {
	PID        int32
	Exec       string
	Status     []string
	CpuPercent float64
	MemUsage   float64
}

func collectProcessById(pid int32) (*ProcessInfo, error) {
	p, _ := process.NewProcess(pid)
	exe, err := p.Exe()
	if err != nil {
		return nil, err
	}

	status, _ := p.Status()
	cpuPercent, _ := p.CPUPercent()
	memStat, _ := p.MemoryInfo()
	processInfo := &ProcessInfo{
		PID:        p.Pid,
		Exec:       exe,
		Status:     status,
		CpuPercent: Round(cpuPercent),
		MemUsage:   Round(ToMB(memStat.RSS)),
	}

	return processInfo, nil

}

func collectProcess(sysInfo *SysInfo) {
	config := config.GetSysInfoReportConfig()
	processInfo, _ := process.Processes()
	sysInfo.Processes = make([]ProcessInfo, len(config.ProcessMonitor))
	for index, _ := range processInfo {

		p, err := collectProcessById(processInfo[index].Pid)
		if err == nil && funk.IndexOf(config.ProcessMonitor, p.Exec) > -1 {
			sysInfo.Processes = append(sysInfo.Processes, p)
		}

	}

}
