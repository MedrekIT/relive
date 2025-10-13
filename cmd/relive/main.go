package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/MedrekIT/relive/internal/runner"
	"github.com/MedrekIT/relive/internal/watcher"
)

func main() {
	if len(os.Args) < 1 || len(os.Args) > 2 {
		log.Fatal("\nUsage:\nrelive <project_path> OR relive (in project's directory)")
	}

	fmt.Println("Relive started")
	projectPath := "."
	var err error
	if len(os.Args) == 2 {
		projectPath, err = filepath.Abs(os.Args[1])
		if err != nil {
			log.Fatal("\nError: invalid file path\n\nUsage:\nrelive <project_path> OR relive (in project's directory)")
		}
	}

	cfg := watcher.Config{
		ProjectPath:    projectPath,
		CheckInterval:  time.Millisecond * 500,
		SearchInterval: time.Second * 5,
	}

	cfg.Cmd, err = runner.RunCommand(cfg.ProjectPath, &exec.Cmd{})
	if err != nil {
		log.Fatalf("\nError: %v\n", err)
	}
	go cfg.InspectLoop(cfg.Cmd)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cfg.Cmd, err = runner.StopCommand(cfg.Cmd)
	if err != nil {
		log.Fatalf("\nError: %v\n", err)
	}
	fmt.Println("Relive finished")
}
