package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/glowfi/ghkd/internal/app"
	"github.com/glowfi/ghkd/internal/cli"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse CLI
	opts, err := cli.Parse()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// Create app config and daemon
	inputDir := "/dev/input"
	pidFilePath := filepath.Join(os.TempDir(), "ghkd.pid")
	appConfig := app.NewConfig(inputDir, opts.ConfigPath, pidFilePath)
	daemon := app.NewDaemon(appConfig)

	// Handle command (version, kill, reload, background)
	shouldExit, err := daemon.HandleCommand(opts.Command)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}
	if shouldExit {
		return
	}

	// Run daemon
	if err := daemon.Run(ctx); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
