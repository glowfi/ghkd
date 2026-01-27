package executor

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/glowfi/ghkd/internal/config"
)

// Executor runs commands and scripts
type Executor struct {
	mu      sync.Mutex
	running map[string]*exec.Cmd
}

// New creates a new executor
func New() *Executor {
	return &Executor{
		running: make(map[string]*exec.Cmd),
	}
}

// Execute runs the action for a keybinding
func (e *Executor) Execute(ctx context.Context, kb *config.Keybinding) error {
	switch {
	case kb.Run != "":
		return e.executeCommand(ctx, kb)
	case kb.Script != "" && kb.Interpreter != "":
		return e.executeScript(ctx, kb)
	case kb.File != "":
		return e.executeFile(ctx, kb)
	default:
		return fmt.Errorf("no action defined for keybinding: %s", kb.Name)
	}
}

// executeCommand runs a simple command
func (e *Executor) executeCommand(ctx context.Context, kb *config.Keybinding) error {
	cmd := exec.CommandContext(ctx, "sh", "-c", kb.Run)

	// Don't wait for command to finish (async)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start command: %w", err)
	}

	// Track running command
	e.trackCommand(kb.Name, cmd)

	// Wait in background and clean up
	go func() {
		if err := cmd.Wait(); err != nil {
			return
		}
		e.untrackCommand(kb.Name)
	}()

	return nil
}

// executeScript runs an inline script with interpreter
func (e *Executor) executeScript(ctx context.Context, kb *config.Keybinding) error {
	// Create temp file
	tmpFile, err := os.CreateTemp("/tmp", fmt.Sprintf("hotkeysd-*%s", kb.Name))
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}

	tmpPath := tmpFile.Name()

	// Write script content
	content := kb.Script
	if _, err := tmpFile.WriteString(content); err != nil {
		if err := tmpFile.Close(); err != nil {
			return err
		}
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
		return fmt.Errorf("write script: %w", err)
	}
	tmpFile.Close()

	// Make executable
	os.Chmod(tmpPath, 0o700)

	// Execute
	cmd := exec.CommandContext(ctx, kb.Interpreter, tmpPath)

	if err := cmd.Start(); err != nil {
		if err := os.Remove(tmpPath); err != nil {
			return err
		}
		return fmt.Errorf("start script: %w", err)
	}

	e.trackCommand(kb.Name, cmd)

	// Wait and clean up
	go func() {
		if err := cmd.Wait(); err != nil {
			return
		}
		if err := os.Remove(tmpPath); err != nil {
			return
		}
		e.untrackCommand(kb.Name)
	}()

	return nil
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path // Return original if home dir fails
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// executeFile runs an external script file
func (e *Executor) executeFile(ctx context.Context, kb *config.Keybinding) error {
	path := expandPath(kb.File)

	// Check file exists
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	var cmd *exec.Cmd

	// Check if executable
	if info.Mode()&0o111 != 0 {
		// Executable - run directly
		cmd = exec.CommandContext(ctx, path)
	} else {
		// Not executable - run with sh
		cmd = exec.CommandContext(ctx, "sh", path)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start file: %w", err)
	}

	e.trackCommand(kb.Name, cmd)

	go func() {
		if err := cmd.Wait(); err != nil {
			return
		}
		e.untrackCommand(kb.Name)
	}()

	return nil
}

// trackCommand adds command to running map
func (e *Executor) trackCommand(name string, cmd *exec.Cmd) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.running[name] = cmd
}

// untrackCommand removes command from running map
func (e *Executor) untrackCommand(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.running, name)
}

// IsRunning checks if a keybinding's command is still running
func (e *Executor) IsRunning(name string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, exists := e.running[name]
	return exists
}

func (e *Executor) Shutdown() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	runningCount := len(e.running)
	if runningCount == 0 {
		return nil
	}

	var errs error
	for _, cmd := range e.running {
		// Graceful termination
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			errs = errors.Join(errs, err)
		}
	}

	return errs
}
