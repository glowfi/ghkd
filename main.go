package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/glowfi/ghkd/internal/config"
	"github.com/glowfi/ghkd/internal/executor"
	"github.com/glowfi/ghkd/internal/listener"
	"github.com/glowfi/ghkd/internal/registry"
)

func main() {
	// Create cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Create executor
	exec := executor.New()

	// Create registry
	reg := registry.NewRegistry(cfg.Keybindings)

	// Start keyboard listener
	lst := listener.NewListener()
	if err := lst.Start(ctx); err != nil {
		log.Fatalf("listener: %v", err)
	}

	fmt.Println("ghkd started. Press Ctrl+C to exit.")
	fmt.Printf("Loaded %d keybindings\n", len(cfg.Keybindings))

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Main event loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-lst.Events():
				if !ok {
					return
				}

				pressed := lst.PressedKeys()

				// Skip if no keys pressed
				if len(pressed) == 0 {
					continue
				}

				if match := reg.Match(pressed); match != nil {
					if err := exec.Execute(match); err != nil {
						fmt.Printf("  Error: %v\n", err)
					}
				}
			}
		}
	}()

	// Wait for shutdown signal
	sig := <-sigChan
	fmt.Printf("\nReceived %s, shutting down...\n", sig)

	// Cancel context first (stops event loop)
	cancel()

	// Stop listener
	lst.Stop()

	// Shutdown executor
	if err := exec.Shutdown(); err != nil {
		log.Printf("Shutdown error: %v", err)
	}

	fmt.Println("Goodbye!")
}
