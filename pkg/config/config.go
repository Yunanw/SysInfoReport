package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type SysInfoReportConfig struct {
	Server         string
	ShowNetwork    bool
	ShowCPU        bool
	ShowHost       bool
	ShowMemory     bool
	ShowPartitions bool
	//ServiceCollector []collector.ServiceCollector
	//ProcessCollector []collector.ProcessCollector
}

func setToViper(config *SysInfoReportConfig) {
	viper.Set("Server", config.Server)
	viper.Set("ShowCPU", config.ShowCPU)
	viper.Set("ShowHost", config.ShowHost)
	viper.Set("ShowMemory", config.ShowMemory)
	viper.Set("ShowPartitions", config.ShowPartitions)
	//viper.Set("ServiceCollector", config.ServiceCollector)
	//viper.Set("ProcessCollector", config.ProcessCollector)
}

func getFromViper() (*SysInfoReportConfig, error) {

	config := &SysInfoReportConfig{
		Server:         viper.GetString("Server"),
		ShowCPU:        viper.GetBool("ShowCPU"),
		ShowHost:       viper.GetBool("ShowHost"),
		ShowMemory:     viper.GetBool("ShowMemory"),
		ShowPartitions: viper.GetBool("ShowPartitions"),
	}
	//serviceMonitor := make([]collector.ServiceCollector, 1)
	//if err := viper.UnmarshalKey("ServiceCollector", &serviceMonitor); err != nil {
	//	return nil, err
	//}
	//config.ServiceCollector = serviceMonitor
	//
	//processMonitor := make([]collector.ProcessCollector, 1)
	//
	//if err := viper.UnmarshalKey("ProcessCollector", &processMonitor); err != nil {
	//	return nil, err
	//}
	//config.ProcessCollector = processMonitor
	return config, nil
}

func DefaultConfig() *SysInfoReportConfig {

	//serviceMonitor := make([]collector.ServiceCollector, 0)
	//serviceMonitor = append(serviceMonitor,collector.ServiceCollector{
	//	Name: "edgeupdate",
	//	ProcessMonitor: &collector.ProcessCollector{
	//		Exec:           "",
	//		ShowCPUPercent: true,
	//		ShowMemory:     true,
	//	},
	//})
	//
	//processMonitor := make([]collector.ProcessCollector, 0)
	//processMonitor = append(processMonitor, collector.ProcessCollector{
	//	Exec:           "C:\\Windows\\explorer.exe",
	//	ShowCPUPercent: true,
	//	ShowMemory:     true,
	//})

	ret := SysInfoReportConfig{
		Server:         ":9090",
		ShowCPU:        true,
		ShowHost:       true,
		ShowMemory:     true,
		ShowPartitions: true,
		//ServiceCollector: serviceMonitor,
		//ProcessCollector: processMonitor,
	}

	return &ret
}

func InitConfig(config *SysInfoReportConfig) error {

	if config == nil {
		config = DefaultConfig()
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("SysInfoReport")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			if saveError := SaveToLocalConfig(config); saveError != nil {
				return errors.Wrap(err, "保存到配置文件失败")
			}

		} else {
			return errors.Wrap(err, "配置文件读取失败")
		}
	}
	return nil
}

func SaveToLocalConfig(config *SysInfoReportConfig) error {
	setToViper(config)
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func GetSysInfoReportConfig() *SysInfoReportConfig {

	var (
		err    error
		config *SysInfoReportConfig
	)

	if config, err = getFromViper(); err != nil {
		log.Fatal(err)
	}
	return config
}
