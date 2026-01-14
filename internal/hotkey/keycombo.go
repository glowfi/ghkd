package hotkey

import (
	"fmt"
	"strings"
)

// KeyCombo represents a parsed key combination
type KeyCombo struct {
	Modifiers []uint16 // Modifier key codes (ctrl, alt, shift, super)
	Key       uint16   // Main key code (non-modifier)
	Raw       string   // Original string (ctr+shift+b)
}

func ParseKeyCombo(s string) (KeyCombo, error) {
	combo := KeyCombo{Raw: s}

	s = strings.TrimSpace(s)
	if s == "" {
		return KeyCombo{}, fmt.Errorf("empty key combination")
	}

	parts := strings.Split(s, "+")
	if len(parts) < 2 {
		return KeyCombo{}, fmt.Errorf("atleast one combination of modifier and non-modifier key must be provided")
	}

	var modifiers []uint16
	var nonModifiers []uint16

	for _, part := range parts {
		code, found := LookupKeyCode(part)
		if !found {
			return KeyCombo{}, fmt.Errorf("unknown key: '%s'", part)
		}

		if IsModifier(code) {
			modifiers = append(modifiers, code)
		} else {
			nonModifiers = append(nonModifiers, code)
		}
	}

	if len(modifiers) < 1 {
		return KeyCombo{}, fmt.Errorf("atleast one modifier key must be provided")
	}

	return combo, nil
}
