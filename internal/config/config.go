package config

import (
	"fmt"
	"strings"
)

type Keybinding struct {
	// Identification
	Name string `yaml:"name"`
	Keys string `yaml:"keys"`

	// Action - one of these must be set
	File string `yaml:"file,omitempty"` // External script: "~/script.sh"

	Run string `yaml:"run,omitempty"` // Simple command: "alacritty"

	Interpreter string `yaml:"interpreter,omitempty"` // Script interpreter: "python3,node,bash"
	Script      string `yaml:"script,omitempty"`      // Script content
}

type Config struct {
	Keybindings []Keybinding `yaml:"keybindings"`
}

func (cfg *Config) Validate() error {
	for _, kb := range cfg.Keybindings {
		// name must be given
		if kb.Name == "" {
			return fmt.Errorf("name must be provided")
		}

		keys := strings.Split(kb.Keys, "+")
		if len(keys) < 2 {
			return fmt.Errorf("key must be provided for %s", kb.Name)
		}

		// atleast one modifier key must be present

		// only one non modifier key must be present

		// only one of the action must be present
	}

	return nil
}
