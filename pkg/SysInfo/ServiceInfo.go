package SysInfo

import (
	"SysInfoReport/pkg/config"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/winservices"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"strconv"
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

func collectService(sysInfo *SysInfo) error {

	sysInfoConfig := config.GetSysInfoReportConfig()

	serviceList, err := winservices.ListServices()
	if err != nil {
		return errors.Wrap(err, "读取服务列表失败")
	}
	serviceListMap := map[string]winservices.Service{}
	From(serviceList).
		SelectT(func(s winservices.Service) KeyValue { return KeyValue{Key: s.Name, Value: s} }).
		ToMap(&serviceListMap)
	sysInfo.Services = make([]ServiceInfo, 0)

	for index := range sysInfoConfig.ServiceMonitor {
		sm, ok := serviceListMap[sysInfoConfig.ServiceMonitor[index].Name]
		if ok {
			s, err := winservices.NewService(sm.Name)
			if err != nil {
				log.Printf("%+v\n", (errors.Wrap(err, "取得服务信息失败:"+sm.Name)))
			} else {
				err = s.GetServiceDetail()
				if err == nil {
					serviceInfo := ServiceInfo{
						Name:   s.Name,
						Config: s.Config,
						State:  s.Status.State,
						Pid:    s.Status.Pid,
					}

					if s.Status.Pid != 0 {
						processInfo, err := collectProcessById(int32(s.Status.Pid))
						if err == nil {
							serviceInfo.ProcessInfo = *processInfo
						} else {
							log.Printf("%+v\n", errors.Errorf("取得进程信息失败:%s", strconv.Itoa(int(s.Status.Pid))))
						}
					}

					sysInfo.Services = append(sysInfo.Services, serviceInfo)
				} else {
					log.Printf("%+v\n", errors.Errorf("取得服务信息失败:%s", s.Name))
				}

			}
		} else {
			serviceInfo := ServiceInfo{
				Name: sysInfoConfig.ServiceMonitor[index].Name,
				DESC: "未找到服务",
			}
			sysInfo.Services = append(sysInfo.Services, serviceInfo)
		}

	}
	return err

}
