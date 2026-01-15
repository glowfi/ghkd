package config

import (
	"github.com/glowfi/ghkd/internal/hotkey"
)

type Keybinding struct {
	// Identification
	Name string          `yaml:"name"`
	Keys hotkey.KeyCombo `yaml:"keys"`

	// Action - one of these must be set
	File string `yaml:"file,omitempty"` // External script: "~/script.sh"

	Run string `yaml:"run,omitempty"` // Simple command: "alacritty"

	Interpreter string `yaml:"interpreter,omitempty"` // Script interpreter: "python3,node,bash"
	Script      string `yaml:"script,omitempty"`      // Script content
}

type Config struct {
	Keybindings []Keybinding `yaml:"keybindings"`
}
