package conf

import (
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/helpers"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func defaultConfigDirectory() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".v2raypanel")
}

func defaultConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), "config.json")
}

func getRuntimeConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), "config.runtime.json")
}

func InitDefaultConfigFile() {
	if err := os.MkdirAll(defaultConfigDirectory(), 0755); err != nil {
		fmt.Println("初始化配置文件夹失败，请手动创建文件夹：", defaultConfigDirectory())
		os.Exit(0)
	}
	readConfigsInit()
}

// ActiveRuntimeConfigFile 写入uuid指定的配置文件.
func ActiveRuntimeConfigFile(uuid string) (string, error) {
	config, err := GetConfigByUUID(uuid)
	if err != nil {
		return "", errors.Wrap(err, "配置文件不存在")
	}
	runPath := getRuntimeConfigPath()

	// 将本地配置写入到 config.
	if err := config.ReadLocalConfig(); err != nil {
		return "", errors.Wrap(err, "读取本地配置失败")
	}

	return runPath, helpers.WriteJSONFile(runPath, config.Config, true)
}

func getLocalConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), "config.local.json")
}

func readConfigsInit() {
	configPath := defaultConfigPath()
	if err := checkFile(configPath, func() error {
		return ioutil.WriteFile(configPath, []byte("{}"), 0755)
	}); err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
		return
	}
}

func checkFile(path string, cb func() error) error {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cb()
		}
		return err
	}
	return nil
}

func GetGfwPath() string {
	return filepath.Join(defaultConfigDirectory(), "gfwlist.txt")
}

func GetActiveRuntimeConfig() string {
	conf, err := getConfigList()
	if err != nil {
		return ""
	}
	for _, item := range conf {
		if item.Status == StatusStart {
			return item.UUID
		}
	}
	return ""
}
