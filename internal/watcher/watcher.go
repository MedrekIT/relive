package watcher

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/MedrekIT/relive/internal/runner"
)

type filesCache struct {
	files map[string]time.Time
	mu    sync.Mutex
}

type Config struct {
	Flags struct {
		Verbose     bool
		ProjectPath string
	}
	Cmd            *exec.Cmd
	ProjectPath    string
	CheckInterval  time.Duration
	SearchInterval time.Duration
}

func (cached *filesCache) searchNewFiles(projectPath string) (bool, error) {
	cached.mu.Lock()
	defer cached.mu.Unlock()

	newFile := false
	err := filepath.Walk(projectPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			if _, ok := cached.files[path]; !ok {
				if err != nil {
					return fmt.Errorf("couldn't access path - %w", err)
				}
				if filepath.Ext(path) == ".go" {
					cached.files[path] = info.ModTime()
					newFile = true
				}
			}
		}
		return nil
	})
	if err != nil {
		return false, fmt.Errorf("couldn't walk through file path - %w", err)
	}

	return newFile, nil
}

func (cached *filesCache) watchForChanges() (bool, error) {
	cached.mu.Lock()
	defer cached.mu.Unlock()

	changes := false
	for file, modTime := range cached.files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			if os.IsNotExist(err) {
				delete(cached.files, file)
			} else {
				return false, fmt.Errorf("couldn't open cached file - %w", err)
			}
		} else {
			if fileInfo.ModTime().After(modTime) {
				cached.files[file] = fileInfo.ModTime()
				changes = true
			}
		}
	}
	return changes, nil
}

func (cfg *Config) InspectLoop(cmd *exec.Cmd) error {
	changesTicker := time.NewTicker(cfg.CheckInterval)
	filesTicker := time.NewTicker(cfg.SearchInterval)
	cachedFiles := filesCache{
		files: make(map[string]time.Time),
	}
	defer func() {
		changesTicker.Stop()
		filesTicker.Stop()
	}()

	changes, err := cachedFiles.searchNewFiles(cfg.ProjectPath)
	if err != nil {
		return err
	}

	for {
		select {
		case <-changesTicker.C:
			changes, err = cachedFiles.watchForChanges()
			if err != nil {
				return err
			}
			if changes {
				if cfg.Flags.Verbose {
					log.Println("Found changes in the project, restarting...")
				}
				cfg.Cmd, err = runner.RerunCommand(cfg.ProjectPath, cfg.Cmd)
				if err != nil {
					return err
				}
			}
		case <-filesTicker.C:
			changes, err = cachedFiles.searchNewFiles(cfg.ProjectPath)
			if err != nil {
				return err
			}
			if changes {
				if cfg.Flags.Verbose {
					log.Println("Found new files in the project, restarting...")
				}
				cfg.Cmd, err = runner.RerunCommand(cfg.ProjectPath, cfg.Cmd)
				if err != nil {
					return err
				}
			}
		}
	}
}
