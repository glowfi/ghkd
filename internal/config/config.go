package config

import (
	"errors"
	"os"

	"github.com/glowfi/ghkd/internal/hotkey"
	"github.com/goccy/go-yaml"
)

var ErrMissingKeybindingName = errors.New("must provide a name to the keybinding")

type Keybinding struct {
	// Identification
	Name           string          `yaml:"name"`
	KeyCombination hotkey.KeyCombo `yaml:"keys"`

	// Action - one of these must be set
	File string `yaml:"file,omitempty"` // External script: "~/script.sh"

	Run string `yaml:"run,omitempty"` // Simple command: "alacritty"

	Interpreter string `yaml:"interpreter,omitempty"` // Script interpreter: "python3,node,bash"
	Script      string `yaml:"script,omitempty"`      // Script content
}

type Config struct {
	Keybindings []Keybinding `yaml:"keybindings"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	for _, kb := range cfg.Keybindings {
		if kb.Name == "" {
			return Config{}, ErrMissingKeybindingName
		}
	}

	return cfg, nil
}
