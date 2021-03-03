package SysInfo

import (
	"SysInfoReport/pkg/config"
	"github.com/shirou/gopsutil/v3/winservices"
	"github.com/thoas/go-funk"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type ServiceInfo struct {
	Name        string
	Config      mgr.Config
	State       svc.State
	Pid         uint32
	ProcessInfo ProcessInfo
	DESC        string
	// contains filtered or unexported fields
}

func collectService(sysInfo *SysInfo) {

	sysInfoConfig := config.GetSysInfoReportConfig()

	serviceList, _ := winservices.ListServices()
	serviceListMap := funk.ToMap(serviceList, "Name").(map[string]winservices.Service)
	sysInfo.Services = make([]ServiceInfo, 0)

	for index, _ := range sysInfoConfig.ServiceMonitor {
		sm, ok := serviceListMap[sysInfoConfig.ServiceMonitor[index].Name]
		if ok {
			s, err := winservices.NewService(sm.Name)
			if err == nil {
				s.GetServiceDetail()
				serviceInfo := ServiceInfo{
					Name:   s.Name,
					Config: s.Config,
					State:  s.Status.State,
					Pid:    s.Status.Pid,
				}

				if s.Status.Pid != 0 {
					processInfo, _ := collectProcessById(int32(s.Status.Pid))
					serviceInfo.ProcessInfo = *processInfo
				}

				sysInfo.Services = append(sysInfo.Services, serviceInfo)

			}
		} else {
			serviceInfo := ServiceInfo{
				Name: sysInfoConfig.ServiceMonitor[index].Name,
				DESC: "未找到服务",
			}
			sysInfo.Services = append(sysInfo.Services, serviceInfo)
		}

	}

}
