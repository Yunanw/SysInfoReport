package SysInfo

import (
	"SysInfoReport/pkg/config"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"github.com/shirou/gopsutil/v3/winservices"
	"github.com/thoas/go-funk"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
)

type Service struct {
	Name          string
	Config        mgr.Config
	State         svc.State
	Accepts       svc.Accepted
	Pid           uint32
	Win32ExitCode uint32
	// contains filtered or unexported fields
}

func collectService(sysInfo *SysInfo) {

	sysInfoConfig := config.GetSysInfoReportConfig()

	serviceList, err := winservices.ListServices()
	serviceMonitorMap := funk.ToMap(sysInfoConfig.ServiceMonitor, "Name").(map[string]config.ServiceMonitor)

	serviceList = funk.Filter(serviceList, func(s winservices.Service) bool {
		_, ok := serviceMonitorMap[s.Name]
		return ok
	}).([]winservices.Service)

	for index, _ := range serviceList {
		s, _ := winservices.NewService(serviceList[index].Name)
		s.GetServiceDetail()
		serviceList[index] = *s

		if s.Status.Pid != 0 {
			proc, _ := process.NewProcess(int32(s.Status.Pid))
			fmt.Println(proc)
		}

	}

	fmt.Println(serviceList)
	fmt.Println(err)
}
