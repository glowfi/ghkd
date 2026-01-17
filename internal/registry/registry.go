package registry

import (
	"sync"

	"github.com/glowfi/ghkd/internal/config"
)

type Registry struct {
	mu       sync.RWMutex // Read-Write Mutex
	bindings []config.Keybinding
}

// NewRegistry creates a new registry
func NewRegistry(bindings []config.Keybinding) *Registry {
	return &Registry{
		bindings: bindings,
	}
}

// Update replaces the current keybindings with new ones (Thread-Safe)
func (r *Registry) Update(bindings []config.Keybinding) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.bindings = bindings
}

// Match finds a keybinding matching pressed keys (Thread-Safe)
func (r *Registry) Match(pressed []uint16) *config.Keybinding {
	r.mu.RLock() // Read lock allows multiple readers, blocks writers
	defer r.mu.RUnlock()

	for i := range r.bindings {
		// Use pointer to avoid copying
		kb := &r.bindings[i]
		if kb.KeyCombination.Matches(pressed) {
			return kb
		}
	}
	return nil
}
