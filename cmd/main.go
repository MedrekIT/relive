package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/MedrekIT/arun/internal/runner"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage:\narun <project_path>")
	}

	project_path, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal("Error: invalid file path\n\nUsage:\narun <project_path>")
	}

	cmd, err := runner.RunCommand(project_path)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	// go watcher.WatchForChanges(time.Second * 30)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("Finishing run...")
	if cmd != nil && cmd.Process != nil {
		cmd.Process.Kill()
		cmd.Wait()
	}
	log.Println("Arun finished")
}
