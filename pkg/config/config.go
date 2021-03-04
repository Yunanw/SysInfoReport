package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

type ServiceMonitor struct {
	Name           string
	ProcessMonitor ProcessMonitor
}

type ProcessMonitor struct {
	Exec           string
	ShowCPUPercent bool
	ShowMemory     bool
}

type SysInfoReportConfig struct {
	Server         string
	ShowNetwork    bool
	ShowCPU        bool
	ShowHost       bool
	ShowMemory     bool
	ShowPartitions bool
	ServiceMonitor []ServiceMonitor
	ProcessMonitor []ProcessMonitor
}

func setToViper(config *SysInfoReportConfig) {
	viper.Set("Server", config.Server)
	viper.Set("ShowCPU", config.ShowCPU)
	viper.Set("ShowHost", config.ShowHost)
	viper.Set("ShowMemory", config.ShowMemory)
	viper.Set("ShowPartitions", config.ShowPartitions)
	viper.Set("ServiceMonitor", config.ServiceMonitor)
	viper.Set("ProcessMonitor", config.ProcessMonitor)
}

func getFromViper() (*SysInfoReportConfig, error) {

	config := &SysInfoReportConfig{
		Server:         viper.GetString("Server"),
		ShowCPU:        viper.GetBool("ShowCPU"),
		ShowHost:       viper.GetBool("ShowHost"),
		ShowMemory:     viper.GetBool("ShowMemory"),
		ShowPartitions: viper.GetBool("ShowPartitions"),
	}
	serviceMonitor := make([]ServiceMonitor, 1)
	if err := viper.UnmarshalKey("ServiceMonitor", &serviceMonitor); err != nil {
		return nil, err
	}
	config.ServiceMonitor = serviceMonitor

	processMonitor := make([]ProcessMonitor, 1)

	if err := viper.UnmarshalKey("ProcessMonitor", &processMonitor); err != nil {
		return nil, err
	}
	config.ProcessMonitor = processMonitor
	return config, nil
}

func DefaultConfig() *SysInfoReportConfig {

	serviceMonitor := make([]ServiceMonitor, 0)
	serviceMonitor = append(serviceMonitor, ServiceMonitor{
		Name: "edgeupdate",
		ProcessMonitor: ProcessMonitor{
			Exec:           "",
			ShowCPUPercent: true,
			ShowMemory:     true,
		},
	})

	processMonitor := make([]ProcessMonitor, 0)
	processMonitor = append(processMonitor, ProcessMonitor{
		Exec:           "C:\\Windows\\explorer.exe",
		ShowCPUPercent: true,
		ShowMemory:     true,
	})

	ret := SysInfoReportConfig{
		Server:         ":9090",
		ShowCPU:        true,
		ShowHost:       true,
		ShowMemory:     true,
		ShowPartitions: true,
		ServiceMonitor: serviceMonitor,
		ProcessMonitor: processMonitor,
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
