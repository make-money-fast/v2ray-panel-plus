package conf

import (
	_ "embed"
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/helpers"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 客户端名称
var (
	// ConfigJsonName 列表
	ConfigJsonName = "config.json"

	// RuntimeConfigJsonName 运行时
	RuntimeConfigJsonName = "config.runtime.json"

	// ConfigLocalName 本地配置
	ConfigLocalName = "config.local.json"

	// ServerConfigTemplate 服务端运行模板文件
	ServerConfigTemplate = "config.template.json"
)

// AsServer 服务端名称
func AsServer() {
	ConfigJsonName = "config.server.json"
	RuntimeConfigJsonName = "config.runtime.json"
}

func GetDefaultConfigDirectory() string {
	return defaultConfigDirectory()
}

func GetRuntimeConfigPath() string {
	return getRuntimeConfigPath()
}

func defaultConfigDirectory() string {
	home := os.Getenv("HOME")
	return filepath.Join(home, ".v2raypanel")
}

func defaultConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), ConfigJsonName)
}

func defaultServerTemplate() string {
	return filepath.Join(defaultConfigDirectory(), ServerConfigTemplate)
}

func getRuntimeConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), RuntimeConfigJsonName)
}

func InitDefaultConfigFile() {
	if err := os.MkdirAll(defaultConfigDirectory(), 0755); err != nil {
		fmt.Println("初始化配置文件夹失败，请手动创建文件夹：", defaultConfigDirectory())
		os.Exit(0)
	}
	readConfigsInit()
}

var (
	//go:embed default_server_config.json
	templateJson string
)

func InitTemplateFile() {
	tplPath := defaultServerTemplate()
	if err := checkFile(tplPath, func() error {
		return ioutil.WriteFile(tplPath, []byte(templateJson), 0755)
	}); err != nil {
		log.Println("初始化配置文件失败", err)
		os.Exit(1)
		return
	}
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

// ActiveServerRuntimeConfig 写入uuid指定的配置文件.
func ActiveServerRuntimeConfig() (string, error) {
	conf, err := MergeServerConfig()
	if err != nil {
		return "", err
	}
	runPath := getRuntimeConfigPath()
	return runPath, helpers.WriteJSONFile(runPath, conf, true)
}

func getLocalConfigPath() string {
	return filepath.Join(defaultConfigDirectory(), ConfigLocalName)
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
