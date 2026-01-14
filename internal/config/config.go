package config

type Config struct {
	Keybindings []Keybinding `yaml:"keybindings"`
}

type Keybinding struct {
	// Identification
	Name string `yaml:"name,omitempty"`
	Keys string `yaml:"keys,omitempty"`

	// Action - one of these must be set
	File string `yaml:"file,omitempty"` // External script: "~/script.sh"

	Run string `yaml:"run,omitempty"` // Simple command: "alacritty"

	Interpreter string `yaml:"interpreter,omitempty"` // Script interpreter: "python3,node,bash"
	Script      string `yaml:"script,omitempty"`      // Script content
}
