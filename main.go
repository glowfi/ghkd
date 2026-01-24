package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/glowfi/ghkd/internal/config"
	"github.com/glowfi/ghkd/internal/daemon"
	"github.com/glowfi/ghkd/internal/executor"
	"github.com/glowfi/ghkd/internal/listener"
	"github.com/glowfi/ghkd/internal/registry"
)

const version = "1.0.0"

func main() {
	// 1. Define Flags
	var (
		configPath  string
		background  bool
		kill        bool
		reload      bool
		showVersion bool
	)

	// Bind both short (-c) and long (--config) to the same variable
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&configPath, "config", "config.yaml", "config path")

	flag.BoolVar(&background, "b", false, "background")
	flag.BoolVar(&background, "background", false, "background")

	flag.BoolVar(&kill, "k", false, "kill")
	flag.BoolVar(&kill, "kill", false, "kill")

	flag.BoolVar(&reload, "r", false, "reload")
	flag.BoolVar(&reload, "reload", false, "reload")

	flag.BoolVar(&showVersion, "v", false, "version")
	flag.BoolVar(&showVersion, "version", false, "version")

	// 2. Custom Help Message
	flag.Usage = func() {
		fmt.Print(`ghkd - Go Hotkey Daemon

Usage:
  ghkd [flags]

Flags:
  -h,  --help              Prints this help message
  -c,  --config [path]     Reads the config from custom path
  -b,  --background        Runs ghkd in the background
  -k,  --kill              Gracefully kills running instances
  -r,  --reload            Reloads configuration of running instance
  -v,  --version           Prints current version
`)
	}

	flag.Parse()

	// 3. Handle Immediate Actions (Version, Kill, Reload)

	if showVersion {
		fmt.Printf("ghkd version %s\n", version)
		return
	}

	if kill {
		if err := daemon.KillInstance(syscall.SIGTERM); err != nil {
			fmt.Printf("Error killing daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sent SIGTERM to ghkd daemon.")
		return
	}

	if reload {
		if err := daemon.KillInstance(syscall.SIGHUP); err != nil {
			fmt.Printf("Error reloading daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sent SIGHUP to ghkd daemon.")
		return
	}

	// 4. Handle Backgrounding
	if background {
		// Re-run the command without the -b flag
		newArgs := []string{}
		for _, arg := range os.Args[1:] {
			if arg != "-b" && arg != "--background" {
				newArgs = append(newArgs, arg)
			}
		}

		cmd := exec.Command(os.Args[0], newArgs...)
		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start background process: %v", err)
		}
		fmt.Printf("ghkd started in background (PID: %d)\n", cmd.Process.Pid)
		return
	}

	// ==========================================
	// 5. Normal Daemon Startup
	// ==========================================

	// 1. CHECK IF ALREADY RUNNING
	if isRunning() {
		fmt.Println("ghkd is already running. Exiting.")
		os.Exit(1)
	}

	// 2. Write PID file
	if err := daemon.WritePID(); err != nil {
		log.Printf("Warning: could not write PID file: %v", err)
	}
	defer daemon.RemovePID()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runDaemon(ctx, configPath)
}

func isRunning() bool {
	pid, err := daemon.ReadPID()
	if err != nil {
		// No PID file found, safe to assume not running
		return false
	}

	// PID file exists, check if the process is actually alive
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists without killing it
	// On Unix, this returns nil if process exists, error if it doesn't
	if err := process.Signal(syscall.Signal(0)); err == nil {
		return true
	}

	// Process listed in PID file is dead, clean up the stale file
	daemon.RemovePID()
	return false
}

func runDaemon(ctx context.Context, configPath string) {
	// Load Config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	exec := executor.New()
	reg := registry.NewRegistry(cfg.Keybindings)

	lst := listener.NewListener()
	if err := lst.Start(ctx); err != nil {
		log.Fatalf("Listener error: %v", err)
	}

	// Handle Signals (SIGINT, SIGTERM, SIGHUP)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Event Loop
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
				if len(pressed) == 0 {
					continue
				}

				if match := reg.Match(pressed); match != nil {
					// Execute in goroutine to not block listener
					go func(m *config.Keybinding) {
						if err := exec.Execute(ctx, m); err != nil {
							fmt.Printf("Error: %v\n", err)
						}
					}(match)
				}
			}
		}
	}()

	// Signal Loop
	for {
		sig := <-sigChan
		if sig == syscall.SIGHUP {
			// Reload Logic
			fmt.Println("Reloading config...")
			newCfg, err := config.LoadConfig(configPath)
			if err != nil {
				log.Printf("Reload failed: %v", err)
				continue
			}
			reg.Update(newCfg.Keybindings)
			fmt.Printf("Reloaded %d keybindings\n", len(newCfg.Keybindings))
			continue
		}

		// Shutdown
		fmt.Println("\nShutting down...")
		break
	}

	// Cleanup
	lst.Stop()
	if err := exec.Shutdown(); err != nil {
		log.Println(err)
	}
}
