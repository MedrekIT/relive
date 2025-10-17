package runner

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func RerunCommand(projectPath string, cmd *exec.Cmd) (*exec.Cmd, error) {
	_, err := StopCommand(cmd)
	if err != nil {
		return &exec.Cmd{}, err
	}
	return RunCommand(projectPath)
}

func StopCommand(cmd *exec.Cmd) (*exec.Cmd, error) {
	if cmd != nil && cmd.Process != nil {
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
			return &exec.Cmd{}, fmt.Errorf("couldn't kill the process - %w", err)
		}
		_ = cmd.Wait()
	}

	return cmd, nil
}

func RunCommand(projectPath string) (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", "./...")
	cmd.Dir = projectPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	err := cmd.Start()
	if err != nil {
		return &exec.Cmd{}, fmt.Errorf("couldn't run a command - %w", err)
	}

	return cmd, nil
}
