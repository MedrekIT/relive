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

var projectPath *string = flag.StringP("project", "p", os.Getenv("PWD"), "Define ReLive entrypoint")
var verbose *bool = flag.BoolP("verbose", "v", false, "Shows ReLive logs")

func main() {
	flag.Parse()

	if len(os.Args) < 1 {
		log.Fatal("\nUsage:\nrelive [OPTIONS] [PATH] OR relive [OPTIONS] (in project's directory)")
	}

	if *verbose {
		fmt.Println("ReLive started")
	}
	var err error
	*projectPath, err = filepath.Abs(*projectPath)
	if err != nil {
		log.Fatal("\nReLive error: invalid file path\n\nUsage:\nrelive [OPTIONS] [PATH] OR relive [OPTIONS] (in project's directory)")
	}

	cfg := watcher.Config{
		Flags: struct {
			Verbose     bool
			ProjectPath string
		}{
			Verbose:     *verbose,
			ProjectPath: *projectPath,
		},
		ProjectPath:    *projectPath,
		CheckInterval:  time.Millisecond * 500,
		SearchInterval: time.Second * 5,
	}

	go func() {
		err = cfg.InspectLoop(cfg.Cmd)
		if err != nil {
			log.Fatalf("\nReLive error: %v\n", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	cfg.Cmd, err = runner.StopCommand(cfg.Cmd)
	if err != nil {
		log.Fatalf("\nReLive error: %v\n", err)
	}
	if *verbose {
		fmt.Println("ReLive finished")
	}
}
