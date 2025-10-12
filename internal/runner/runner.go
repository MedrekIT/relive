package runner

import (
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(project_path string) (*exec.Cmd, error) {
	cmd := exec.Command("go", "run", "./...")
	cmd.Dir = project_path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return &exec.Cmd{}, fmt.Errorf("couldn't run a command - %w", err)
	}

	return cmd, nil
}
