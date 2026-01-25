package pid

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

type PidManager struct {
	pidFile string
}

func NewPidManager(pidFile string) *PidManager {
	return &PidManager{
		pidFile: pidFile,
	}
}

// WritePID writes the current PID to a file
func (p *PidManager) WritePID() error {
	pid := os.Getpid()
	return os.WriteFile(p.pidFile, []byte(strconv.Itoa(pid)), 0o644)
}

// RemovePID removes the PID file
func (p *PidManager) RemovePID() {
	os.Remove(p.pidFile)
}

// ReadPID reads the PID from the file
func (p *PidManager) ReadPID() (int, error) {
	data, err := os.ReadFile(p.pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

// KillInstance sends a signal to the running daemon
func (p *PidManager) KillInstance(sig syscall.Signal) error {
	pid, err := p.ReadPID()
	if err != nil {
		return fmt.Errorf("daemon not running (pid file not found)")
	}

	// Find process and send signal
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}

func (p *PidManager) IsRunning() bool {
	pid, err := p.ReadPID()
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
	p.RemovePID()
	return false
}
