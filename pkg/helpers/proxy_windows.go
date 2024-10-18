//go:build windows

package helpers

import (
	"github.com/pkg/errors"
	"golang.org/x/sys/windows/registry"
)

const (
	internetSetting = `Software\Microsoft\Windows\CurrentVersion\Internet Settings`
)

func SetProxy(addr string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSetting, registry.ALL_ACCESS)
	if err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	}
	defer k.Close()

	err = k.SetStringValue("AutoConfigURL", addr)
	if err != nil {
		return errors.Wrap(err, "设置系统代理失败, 请将系统代理手动设置为: "+addr)
	}
	return nil
}

func UnSetProxy() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, internetSetting, registry.ALL_ACCESS)
	if err != nil {
		return errors.Wrap(err, "清楚系统代理失败, 请手动操作")
	}
	defer k.Close()

	err = k.DeleteValue("AutoConfigURL")
	if err != nil {
		return errors.Wrap(err, "清除系统代理失败, 请手动清除")
	}
	return nil
}
