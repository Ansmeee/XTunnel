package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type TunnelConfig struct {
	FileName   string `json:"file_name"`
	ConfigName string `json:"config_name"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort string `json:"remote_port"`
	SSHIp      string `json:"ssh_ip"`
	SSHPort    string `json:"ssh_port"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	Switch     bool   `json:"switch"`
}

var configPath = "config"

func (c *TunnelConfig) SaveConfigFile() error {
	fileName := fmt.Sprintf("%s/%d.json", configPath, time.Now().UnixMicro())
	c.FileName = fileName
	c.ConfigName = fmt.Sprintf("%s:%s", c.RemoteIP, c.RemotePort)
	fileContent, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, fileContent, 0644); err != nil {
		return err
	}

	log.Printf("save config file success: %s", fileName)
	return nil
}

func (c *TunnelConfig) LoadConfigFile() ([]*TunnelConfig, error) {
	files, err := os.ReadDir(configPath)
	if err != nil {
		log.Println("load config files error:", err.Error())
		return nil, err
	}

	configs := make([]*TunnelConfig, 0)
	for _, file := range files {
		filePath := configPath + "/" + file.Name()
		config, err := os.ReadFile(filePath)
		if err != nil {
			log.Println("read config file error:", file.Name(), err.Error())
			continue
		}

		conf := &TunnelConfig{}
		if err = json.Unmarshal(config, conf); err != nil {
			log.Println("unmarshal config file error:", file.Name(), err.Error())
			continue
		}
		configs = append(configs, conf)
	}

	return configs, nil
}
