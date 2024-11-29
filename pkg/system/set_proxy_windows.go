package system

import (
	"os"
	"os/exec"
)

func SetProxy(addr string) error {
	return nil
}

func UnSetProxy() error {
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
