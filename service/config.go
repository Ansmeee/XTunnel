package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"os"
	"time"
	"xtunnel/logger"
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

func (c *ConfigFile) DeleteConfigFile(ctx context.Context) error {
	fileName := c.FileName
	if fileName == "" {
		return fmt.Errorf("invalid config file")
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		logger.Error(ctx, "config file not exists", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file not exists")
	}

	if err := os.Remove(fileName); err != nil {
		logger.Error(ctx, "config file delete error", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file delete error")
	}

	return nil
}

func (c *ConfigFile) UpdateConfigFile(ctx context.Context) error {
	fileName := c.FileName
	if fileName == "" {
		return fmt.Errorf("invalid config file")
	}

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		logger.Error(ctx, "config file not exists", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file not exists")
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		logger.Error(ctx, "config file open error", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file open error")
	}
	defer file.Close()

	newContent, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = file.Write(newContent)
	if err != nil {
		logger.Error(ctx, "config file write error", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file write error")
	}

	logger.Info(ctx, "config file updated", g.Map{"filename": fileName})
	return nil
}

func (c *ConfigFile) SaveConfigFile(ctx context.Context) error {
	if err := c.EnsureDir(ctx); err != nil {
		return err
	}

	fileName := fmt.Sprintf("%s/%d.json", configPath, time.Now().UnixMicro())
	c.FileName = fileName
	if c.ConfigName == "" {
		c.ConfigName = fmt.Sprintf("%s:%s", c.RemoteIP, c.RemotePort)
	}

	fileContent, err := json.Marshal(c)
	if err != nil {
		logger.Error(ctx, "config file marshal error", g.Map{"filename": fileName, "error": err.Error()})
		return err
	}

	if err := os.WriteFile(fileName, fileContent, 0644); err != nil {
		logger.Error(ctx, "config file write error", g.Map{"filename": fileName, "error": err.Error()})
		return fmt.Errorf("config file write error")
	}

	logger.Info(ctx, "config file saved", g.Map{"filename": fileName})
	return nil
}

func (c *ConfigFile) LoadConfigFile(ctx context.Context) ([]*ConfigFile, error) {
	configs := make([]*ConfigFile, 0)

	if err := c.EnsureDir(ctx); err != nil {
		return configs, nil
	}

	files, err := os.ReadDir(configPath)
	if err != nil {
		logger.Error(ctx, "load config files error", g.Map{"config_path": configPath, "error": err.Error()})
		return nil, fmt.Errorf("load config files error")
	}

	for _, file := range files {
		filePath := configPath + "/" + file.Name()
		config, err := os.ReadFile(filePath)
		if err != nil {
			logger.Error(ctx, "read config file error", g.Map{"file_path": filePath, "error": err.Error()})
			continue
		}

		conf := &ConfigFile{}
		if err = json.Unmarshal(config, conf); err != nil {
			logger.Error(ctx, "unmarshal config file error", g.Map{"file_name": file.Name(), "error": err.Error()})
			continue
		}
		configs = append(configs, conf)
	}

	return configs, nil
}

func (c *ConfigFile) EnsureDir(ctx context.Context) error {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			if err := os.Mkdir(configPath, 0644); err != nil {
				logger.Error(ctx, "config file mkdir error", g.Map{"config_path": configPath, "error": err.Error()})
				return fmt.Errorf("config file mkdir error")
			}
			return nil
		}

		logger.Error(ctx, "config dir is not ready", g.Map{"config_path": configPath, "error": err.Error()})
		return fmt.Errorf("config dir not ready")
	}

	return nil
}
