package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"

	"github.com/MedrekIT/relive/internal/runner"
	"github.com/MedrekIT/relive/internal/watcher"

	flag "github.com/spf13/pflag"
)

var verbose *bool = flag.BoolP("verbose", "v", false, "Shows ReLive logs")

func main() {
	flag.Parse()
	if len(os.Args) < 1 {
		log.Fatal("\nUsage:\nrelive <project_path> [OPTIONS] OR relive [OPTIONS] (in project's directory)")
	}

	if *verbose {
		fmt.Println("ReLive started")
	}
	projectPath := "."
	var err error
	if len(flag.Args()) == 1 {
		projectPath, err = filepath.Abs(flag.Arg(0))
		if err != nil {
			log.Fatal("\nError: invalid file path\n\nUsage:\nrelive <project_path> OR relive (in project's directory)")
		}
	}

	cfg := watcher.Config{
		ProjectPath:    projectPath,
		CheckInterval:  time.Millisecond * 500,
		SearchInterval: time.Second * 5,
	}

	cfg.Cmd, err = runner.RunCommand(cfg.ProjectPath)
	if err != nil {
		log.Fatalf("\nError: %v\n", err)
	}
	go func() {
		err = cfg.InspectLoop(cfg.Cmd)
		if err != nil {
			log.Fatalf("\nError: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cfg.Cmd, err = runner.StopCommand(cfg.Cmd)
	if err != nil {
		log.Fatalf("\nError: %v\n", err)
	}
	if *verbose {
		fmt.Println("ReLive finished")
	}
}
