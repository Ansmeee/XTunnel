package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type ConfigFile struct {
	Identifier string `json:"identifier"`
	FileName   string `json:"file_name"`
	ConfigName string `json:"config_name"`
	RemoteIP   string `json:"remote_ip"`
	RemotePort string `json:"remote_port"`
	ServerIP   string `json:"server_ip"`
	ServerPort string `json:"server_port"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
}

var configPath = "config"

func (c *ConfigFile) DeleteConfigFile() error {
	fileName := c.FileName
	if fileName == "" {
		return fmt.Errorf("invalid config file")
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("[%s] config file not exists", fileName)
		return fmt.Errorf("config file not exists")
	}

	if err := os.Remove(fileName); err != nil {
		log.Printf("[%s] config file delete error: %s", fileName, err.Error())
		return fmt.Errorf("config file delete error")
	}

	return nil
}

func (c *ConfigFile) UpdateConfigFile() error {
	fileName := c.FileName
	if fileName == "" {
		return fmt.Errorf("invalid config file")
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Printf("[%s] config file not exists", fileName)
		return fmt.Errorf("config file not exists")
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("[%s] config file open error: %s", fileName, err.Error())
		return fmt.Errorf("config file open error")
	}
	defer file.Close()

	newContent, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = file.Write(newContent)
	if err != nil {
		log.Printf("[%s] config file write error: %s", fileName, err.Error())
		return fmt.Errorf("config file write error")
	}

	log.Printf("[%s] config file updated", fileName)
	return nil
}

func (c *ConfigFile) SaveConfigFile() error {
	if err := c.EnsureDir(); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s/%d.json", configPath, time.Now().UnixMicro())
	c.FileName = fileName
	if c.ConfigName == "" {
		c.ConfigName = fmt.Sprintf("%s:%s", c.RemoteIP, c.RemotePort)
	}

	fileContent, err := json.Marshal(c)
	if err != nil {
		log.Printf("[%s] config file marshal error: %s", fileName, err.Error())
		return err
	}

	if err := os.WriteFile(fileName, fileContent, 0644); err != nil {
		log.Printf("[%s] config file write error: %s", fileName, err.Error())
		return fmt.Errorf("config file write error")
	}

	log.Printf("[%s] config file saved", fileName)
	return nil
}

func (c *ConfigFile) LoadConfigFile() ([]*ConfigFile, error) {
	configs := make([]*ConfigFile, 0)

	if err := c.EnsureDir(); err != nil {
		return configs, nil
	}

	files, err := os.ReadDir(configPath)
	if err != nil {
		log.Printf("load config files error: %s", err.Error())
		return nil, fmt.Errorf("load config files error")
	}

	for _, file := range files {
		filePath := configPath + "/" + file.Name()
		config, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("[%s] read config file error: %s", file.Name(), err.Error())
			continue
		}

		conf := &ConfigFile{}
		if err = json.Unmarshal(config, conf); err != nil {
			log.Printf("[%s] unmarshal config file error: %s", file.Name(), err.Error())
			continue
		}
		configs = append(configs, conf)
	}

	return configs, nil
}

func (c *ConfigFile) EnsureDir() error {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(configPath, 0644); err != nil {
				log.Printf("[%s] config file mkdir error: %s", configPath, err.Error())
				return fmt.Errorf("config file mkdir error")
			}
			return nil
		}

		fmt.Printf("[%s] config dir not ready", configPath)
		return fmt.Errorf("config dir not ready")
	}

	return nil
}
