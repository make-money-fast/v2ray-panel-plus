package conf

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"v2ray-panel-plus/pkg/helpers"
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
		return ioutil.WriteFile(configPath, []byte("{}"), 0644)
	}); err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
		return
	}
}

func checkFile(path string, cb func() error) error {
	_, err := os.Stat(path)
	if err != nil {
		if err == os.ErrNotExist {
			return cb()
		}
		return err
	}
	return nil
}
