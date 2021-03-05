package collector

import (
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/winservices"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"strconv"
)

type ServiceCollector struct {
	Target []ServiceTarget
}

type ServiceTarget struct {
	Name           string
	ShowCPUPercent bool
	ShowMemory     bool
}

type ServiceInfo struct {
	Name        string
	Config      mgr.Config
	State       svc.State
	Pid         uint32
	ProcessInfo ProcessInfo
	DESC        string
}

func (collector *ServiceCollector) Collect(sysInfo *SysInfo) error {

	for index := range collector.Target {
		svc, err := winservices.NewService(collector.Target[index].Name)
		err = svc.GetServiceDetail()
		if err == nil {
			serviceInfo := ServiceInfo{
				Name:   svc.Name,
				Config: svc.Config,
				State:  svc.Status.State,
				Pid:    svc.Status.Pid,
			}
			if svc.Status.Pid != 0 {
				processInfo, err := collectProcessById(int32(svc.Status.Pid))
				if err == nil {
					serviceInfo.ProcessInfo = *processInfo
				} else {
					err = errors.Errorf("取得进程信息失败:%s", strconv.Itoa(int(svc.Status.Pid)))
					log.Printf("%+v\n", err)
				}
			}
			sysInfo.Services = append(sysInfo.Services, serviceInfo)
		} else {
			serviceInfo := ServiceInfo{
				Name: collector.Target[index].Name,
				DESC: "未找到服务或取得服务信息失败",
			}
			log.Printf("%+v\n", errors.Errorf("取得服务信息失败:%s", collector.Target[index].Name))
			sysInfo.Services = append(sysInfo.Services, serviceInfo)
		}
	}
	return nil

}
