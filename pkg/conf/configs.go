package conf

import (
	"fmt"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/helpers"
	"github.com/make-money-fast/v2ray-panel-plus/pkg/pac"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func CreateOneConfig(config *ClientConfig) error {
	list, err := getConfigList()
	if err != nil {
		return err
	}
	list = append(list, config)
	return saveConfigList(list)
}

func CreateOneServerConfig(config *ServerConfig) error {
	list, err := GetServerConfigList()
	if err != nil {
		return err
	}
	list = append(list, config)
	return saveServerConfigList(list)
}

func UpdateOneConfig(config *ClientConfig) error {
	list, err := getConfigList()
	if err != nil {
		return err
	}
	lo.ForEach(list, func(item *ClientConfig, index int) {
		if item.UUID == config.UUID {
			list[index] = config
		}
	})
	return saveConfigList(list)
}

func UpdateOneServerConfig(config *ServerConfig) error {
	list, err := GetServerConfigList()
	if err != nil {
		return err
	}
	lo.ForEach(list, func(item *ServerConfig, index int) {
		if item.UUID == config.UUID {
			list[index] = config
		}
	})
	return saveServerConfigList(list)
}

func GetConfigList() ([]*ClientConfig, error) {
	return getConfigList()
}

func GetServerConfigList() ([]*ServerConfig, error) {
	return getServerConfig()
}

func GetConfigByUUID(uuid string) (*ClientConfig, error) {
	list, err := getConfigList()
	if err != nil {
		return nil, err
	}
	for _, item := range list {
		if item.UUID == uuid {
			return item, nil
		}
	}
	return nil, errors.New("配置文件不存在")
}

func DeleteOneConfig(uuid string) error {
	list, err := getConfigList()
	if err != nil {
		return err
	}
	newConf := lo.Filter(list, func(item *ClientConfig, index int) bool {
		if item.UUID == uuid {
			return false
		}
		return true
	})
	return saveConfigList(newConf)
}

func DeleteOneServerConfig(uuid string) error {
	list, err := GetServerConfigList()
	if err != nil {
		return err
	}
	newConf := lo.Filter(list, func(item *ServerConfig, index int) bool {
		if item.UUID == uuid {
			return false
		}
		return true
	})
	return saveServerConfigList(newConf)
}

func getConfigList() ([]*ClientConfig, error) {
	var (
		confs     []*ClientConfig
		configMap = make(map[string]*ClientConfig)
	)
	err := helpers.ReadJSONFile(defaultConfigPath(), &configMap)
	if err != nil {
		return nil, err
	}
	for _, item := range configMap {
		confs = append(confs, item)
	}
	sort.Slice(confs, func(i, j int) bool {
		return confs[i].Ts < confs[j].Ts
	})
	return confs, nil
}

func getServerConfig() ([]*ServerConfig, error) {
	var (
		confs     []*ServerConfig
		configMap = make(map[string]*ServerConfig)
	)
	err := helpers.ReadJSONFile(defaultConfigPath(), &configMap)
	if err != nil {
		fmt.Println(defaultConfigPath())
		return nil, err
	}
	for _, item := range configMap {
		confs = append(confs, item)
	}
	sort.Slice(confs, func(i, j int) bool {
		return confs[i].Ts < confs[j].Ts
	})
	return confs, nil
}

func saveConfigList(list []*ClientConfig) error {
	var configMap = make(map[string]*ClientConfig)
	lo.ForEach(list, func(item *ClientConfig, index int) {
		configMap[item.UUID] = item
	})
	return helpers.WriteJSONFile(defaultConfigPath(), configMap)
}

func saveServerConfigList(list []*ServerConfig) error {
	var configMap = make(map[string]*ServerConfig)
	lo.ForEach(list, func(item *ServerConfig, index int) {
		configMap[item.UUID] = item
	})
	return helpers.WriteJSONFile(defaultConfigPath(), configMap)
}

func SaveServerConfigList(list []*ServerConfig) error {
	return saveServerConfigList(list)
}

func SaveConfigList(list []*ClientConfig) error {
	return saveConfigList(list)
}

func GetActiveConfig() (*ClientConfig, error) {
	var v2rayConf ClientConfigV2ray
	if err := helpers.ReadJSONFile(getRuntimeConfigPath(), &v2rayConf); err != nil {
		return nil, errors.Wrap(err, "获取运行配置失败")
	}
	configList, err := getConfigList()
	if err != nil {
		return nil, errors.Wrap(err, "获取配置列表失败")
	}
	for _, item := range configList {
		if item.UUID == v2rayConf.UUID {
			return item, nil
		}
	}
	return nil, errors.New("没有正在运行的配置")
}

func ParsePacJS() string {
	local, err := GetLocalConfig()
	if err != nil {
		return ""
	}
	data, err := ioutil.ReadFile(GetGfwPath())
	if err != nil {
		return ""
	}
	p := pac.ParseGFW(data)
	pacJS := p.ToPacjs("PROXY " + strings.TrimPrefix(local.HttpProxy(), "http://"))
	return pacJS
}

type RunningStatus struct {
	ProxyMode   int    `json:"proxy_mode"`
	ProxyStatus int    `json:"proxy_status"`
	RunningUUID string `json:"running_uuid"`
}

func GetStatusPath() string {
	return filepath.Join(defaultConfigDirectory(), "status.json")
}

func InitRunningStatus() {
	var runtimeConfig RunningStatus
	_, err := os.Stat(GetStatusPath())
	if err != nil {
		if os.IsNotExist(err) {
			helpers.WriteJSONFile(GetStatusPath(), runtimeConfig)
		}
		return
	}
}

func SetRunningStatus(s *RunningStatus) error {
	return helpers.WriteJSONFile(GetStatusPath(), s)
}

func GetRunningStatus() (*RunningStatus, error) {
	var s RunningStatus
	err := helpers.ReadJSONFile(GetStatusPath(), &s)
	if err != nil {
		return &RunningStatus{}, nil
	}
	return &s, nil
}
