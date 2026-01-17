package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

var pidFile = filepath.Join(os.TempDir(), "ghkd.pid")

// WritePID writes the current PID to a file
func WritePID() error {
	pid := os.Getpid()
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0o644)
}

// RemovePID removes the PID file
func RemovePID() {
	os.Remove(pidFile)
}

// ReadPID reads the PID from the file
func ReadPID() (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

// KillInstance sends a signal to the running daemon
func KillInstance(sig syscall.Signal) error {
	pid, err := ReadPID()
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
