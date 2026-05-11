//go:build !windows

package core

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func runDaemon() error {
	if os.Getppid() != 1 {
		args := os.Args
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = nil
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setsid: true,
		}
		if err := cmd.Start(); err != nil {
			return err
		}
		if pidFile != "" {
			os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644)
		}
		logSuccess("Daemon started with PID: %d", cmd.Process.Pid)
		os.Exit(0)
	}
	return nil
}
