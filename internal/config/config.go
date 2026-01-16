package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/glowfi/ghkd/internal/hotkey"
	"github.com/goccy/go-yaml"
)

var (
	ErrMissingKeybindingName   = errors.New("must provide a name to the keybinding")
	ErrMultipleActions         = errors.New("only one of 'run', 'script', 'file' allowed")
	ErrNoAction                = errors.New("must provide one of one of 'run', 'script', 'file'")
	ErrScriptNeedsInterpreter  = errors.New("'script' requires 'interpreter'")
	ErrDuplicateKeybinding     = errors.New("duplicate keybinding found")
	ErrDuplicateKeybindingName = errors.New("duplicate keybinding name found")
)

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

	seenKeybindings := map[string]bool{}
	seenKeybindingsName := map[string]bool{}

	for _, kb := range cfg.Keybindings {
		if kb.Name == "" {
			return Config{}, ErrMissingKeybindingName
		}

		if countActions(kb) == 0 {
			return Config{}, fmt.Errorf("%s: %w", kb.Name, ErrNoAction)
		}

		if countActions(kb) > 1 {
			return Config{}, fmt.Errorf("%s: %w", kb.Name, ErrMultipleActions)
		}

		if kb.Script != "" && kb.Interpreter == "" {
			return Config{}, fmt.Errorf("%s: %w", kb.Name, ErrScriptNeedsInterpreter)
		}

		_, KeyBindingexists := seenKeybindings[kb.KeyCombination.Raw]
		if KeyBindingexists {
			return Config{}, fmt.Errorf("%s: %w", kb.Name, ErrDuplicateKeybinding)
		}

		_, KeyBindingNameexists := seenKeybindingsName[kb.Name]
		if KeyBindingNameexists {
			return Config{}, fmt.Errorf("%s: %w", kb.Name, ErrDuplicateKeybindingName)
		}

		seenKeybindingsName[kb.Name] = true
		seenKeybindings[kb.KeyCombination.Raw] = true
	}

	return cfg, nil
}

func countActions(kb Keybinding) int {
	count := 0
	if kb.Run != "" {
		count++
	}
	if kb.Script != "" {
		count++
	}
	if kb.File != "" {
		count++
	}
	return count
}
