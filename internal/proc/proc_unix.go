//go:build !windows
// +build !windows

package proc

import (
	"os"
	"os/exec"
	"syscall"
)

func setupSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func forwardSignal(cmd *exec.Cmd, sig os.Signal) {
	syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
}
