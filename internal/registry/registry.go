package registry

import (
	"github.com/glowfi/ghkd/internal/config"
)

type Registry struct {
	bindings []config.Keybinding
}

func NewRegistry(keyBindings []config.Keybinding) *Registry {
	return &Registry{
		bindings: keyBindings,
	}
}

// Match finds a keybinding that matches pressed keys
func (r *Registry) Match(pressed []uint16) *config.Keybinding {
	for i := range r.bindings {
		if r.bindings[i].KeyCombination.Matches(pressed) {
			return &r.bindings[i]
		}
	}
	return nil
}
