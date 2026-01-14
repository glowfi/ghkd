package config

type Config struct {
	Keybindings []Keybinding `yaml:"keybindings"`
}

type Keybinding struct {
	// Identification
	Name string `yaml:"name"`
	Keys string `yaml:"keys"`

	// Action - one of these must be set
	File string `yaml:"file"` // External script: "~/script.sh"

	Run string `yaml:"run"` // Simple command: "alacritty"

	Interpreter string `yaml:"interpreter"` // Script interpreter: "python3,node,bash"
	Script      string `yaml:"script"`      // Script content
}
