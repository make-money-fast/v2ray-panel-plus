package helpers

import (
	"github.com/pkg/errors"
	"os"
	"os/exec"
)

func SetProxy(addr string) error {
	/**
	networksetup -getwebproxy <networkservice>
	networksetup -setwebproxy <networkservice> <domain> <port number> <authenticated> <username> <password>
	networksetup -setwebproxystate <networkservice> <on off>
	networksetup -getsecurewebproxy <networkservice>
	networksetup -setsecurewebproxy <networkservice> <domain> <port number> <authenticated> <username> <password>
	networksetup -setsecurewebproxystate <networkservice> <on off>
	*/
	c := exec.Command("networksetup", "setautoproxyurl", "Wi-fi", addr)
	if err := runCmd(c); err != nil {
		return errors.Wrap(err, "设置系统代理失败")
	}
	return nil
}

func UnSetProxy() error {
	/**
	https://stackoverflow.com/questions/26992886/set-proxy-through-windows-command-line-including-login-parameters
	networksetup -setwebproxystate <networkservice> <on off>
	networksetup -setsecurewebproxystate <networkservice> <on off>
	networksetup -setautoproxystate <service line Wi-Fi> off
	*/
	c := exec.Command("networksetup", "setautoproxystate", "Wi-fi", "off")
	if err := runCmd(c); err != nil {
		return errors.Wrap(err, "设置系统代理失败")
	}
	return nil
}

func runCmd(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
