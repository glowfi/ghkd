package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/glowfi/ghkd/internal/cli"
	"github.com/glowfi/ghkd/internal/config"
	"github.com/glowfi/ghkd/internal/executor"
	"github.com/glowfi/ghkd/internal/listener"
	"github.com/glowfi/ghkd/internal/pid"
	"github.com/glowfi/ghkd/internal/registry"
)

type Daemon struct {
	config     *Config
	pidManager *pid.PidManager
}

func NewDaemon(cfg *Config) *Daemon {
	return &Daemon{
		config:     cfg,
		pidManager: pid.NewPidManager(cfg.PidFilePath),
	}
}

// HandleCommand processes the CLI command and returns true if the program should exit
func (d *Daemon) HandleCommand(cmd cli.Command) (shouldExit bool, err error) {
	switch cmd {
	case cli.CommandVersion:
		fmt.Printf("ghkd version %s\n", Version)
		return true, nil

	case cli.CommandKill:
		if err := d.pidManager.KillInstance(syscall.SIGTERM); err != nil {
			return true, fmt.Errorf("error killing daemon: %w", err)
		}
		fmt.Println("Sent SIGTERM to ghkd daemon.")
		return true, nil

	case cli.CommandReload:
		if err := d.pidManager.KillInstance(syscall.SIGHUP); err != nil {
			return true, fmt.Errorf("error reloading daemon: %w", err)
		}
		fmt.Println("Sent SIGHUP to ghkd daemon.")
		return true, nil

	case cli.CommandBackground:
		if err := d.startBackground(); err != nil {
			return true, err
		}
		return true, nil

	default:
		return false, nil
	}
}

func (d *Daemon) startBackground() error {
	newArgs := cli.FilterBackgroundFlag(os.Args[1:])
	cmd := exec.Command(os.Args[0], newArgs...)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start background process: %w", err)
	}
	fmt.Printf("ghkd started in background (PID: %d)\n", cmd.Process.Pid)
	return nil
}

func (d *Daemon) Run(ctx context.Context) error {
	// Check if already running
	if d.pidManager.IsRunning() {
		return fmt.Errorf("ghkd is already running")
	}

	// Write PID file
	if err := d.pidManager.WritePID(); err != nil {
		log.Printf("Warning: could not write PID file: %v", err)
	}
	defer d.pidManager.RemovePID()

	return d.runEventLoop(ctx)
}

func (d *Daemon) runEventLoop(ctx context.Context) error {
	// Load Config
	cfg, err := config.LoadConfig(d.config.CfgPath)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	exec := executor.New()
	reg := registry.NewRegistry(cfg.Keybindings)

	lst := listener.NewListener(d.config.InputDir)
	if err := lst.Start(ctx); err != nil {
		return fmt.Errorf("listener error: %w", err)
	}

	// Handle Signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// Event Loop
	go d.processEvents(ctx, lst, reg, exec)

	// Signal Loop
	d.handleSignals(sigChan, reg)

	// Cleanup
	lst.Stop()
	if err := exec.Shutdown(); err != nil {
		log.Println(err)
	}

	return nil
}

func (d *Daemon) processEvents(ctx context.Context, lst *listener.Listener, reg *registry.Registry, exec *executor.Executor) {
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
				go func(m *config.Keybinding) {
					if err := exec.Execute(ctx, m); err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				}(match)
			}
		}
	}
}

func (d *Daemon) handleSignals(sigChan <-chan os.Signal, reg *registry.Registry) {
	for sig := range sigChan {
		if sig == syscall.SIGHUP {
			fmt.Println("Reloading config...")
			newCfg, err := config.LoadConfig(d.config.CfgPath)
			if err != nil {
				log.Printf("Reload failed: %v", err)
				continue
			}
			reg.Update(newCfg.Keybindings)
			fmt.Printf("Reloaded %d keybindings\n", len(newCfg.Keybindings))
			continue
		}

		fmt.Println("\nShutting down...")
		return
	}
}
