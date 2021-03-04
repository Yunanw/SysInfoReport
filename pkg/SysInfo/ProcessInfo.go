package SysInfo

import (
	"SysInfoReport/pkg/config"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/process"
	"log"
	"strconv"
)

type ProcessInfo struct {
	PID        int32
	Exec       string
	CpuPercent float64
	MemUsage   float64
	DESC       string
}

func collectProcessById(pid int32) (*ProcessInfo, error) {
	p, _ := process.NewProcess(pid)
	exe, err := p.Exe()
	if err != nil {
		return nil, err
	}

	cpuPercent, _ := p.CPUPercent()
	memStat, _ := p.MemoryInfo()
	processInfo := &ProcessInfo{
		PID:        p.Pid,
		Exec:       exe,
		CpuPercent: Round(cpuPercent),
		MemUsage:   Round(ToMB(memStat.RSS)),
	}

	return processInfo, nil

}

func collectProcess(sysInfo *SysInfo) error {
	reportConfig := config.GetSysInfoReportConfig()
	processInfo, err := process.Processes()
	if err != nil {
		return errors.Wrap(err, "读取进程列表失败")
	}
	processInfoMap := makeProcessListToExeMap(processInfo)
	for index := range reportConfig.ProcessMonitor {

		psInfo, ok := processInfoMap[reportConfig.ProcessMonitor[index].Exec]
		if ok {
			p, err := collectProcessById(psInfo.PID)
			if err == nil {
				psInfo = *p
			} else {
				log.Printf("%+v\n", errors.Errorf("取得进程信息失败:%s", strconv.Itoa(int(psInfo.PID))))
			}

		} else {
			psInfo.DESC = "未找到进程"
			psInfo.Exec = reportConfig.ProcessMonitor[index].Exec
		}
		sysInfo.Processes = append(sysInfo.Processes, psInfo)
	}

	return err

}

func makeProcessListToExeMap(processInfo []*process.Process) map[string]ProcessInfo {
	processInfoMap := make(map[string]ProcessInfo)
	for index := range processInfo {
		if processInfo[index].Pid == 0 {
			continue
		}
		p, err := fetchProcessExec(processInfo[index].Pid)
		if err == nil {
			processInfoMap[p.Exec] = *p
		}
	}
	return processInfoMap
}

func fetchProcessExec(pid int32) (*ProcessInfo, error) {
	p, _ := process.NewProcess(pid)
	exe, err := p.Exe()
	if err != nil {
		return nil, err
	}
	processInfo := &ProcessInfo{
		PID:  p.Pid,
		Exec: exe,
	}

	return processInfo, nil
}
