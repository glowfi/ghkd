package hotkey

import (
	"errors"
	"strings"
)

var (
	ErrInvalidKeyComboFormat   = errors.New("atleast one combination of modifier and non-modifier key must be provided")
	ErrInvalidNonModifierCount = errors.New("exactly one non‑modifier key must be provided")
	ErrUnknownKey              = errors.New("unknown key")
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
		return KeyCombo{}, ErrInvalidKeyComboFormat
	}

	parts := strings.Split(s, "+")
	if len(parts) < 2 {
		return KeyCombo{}, ErrInvalidKeyComboFormat
	}

	var modifiers []uint16
	var nonModifiers []uint16

	for _, part := range parts {
		code, found := LookupKeyCode(part)
		if !found {
			return KeyCombo{}, ErrUnknownKey
		}

		if IsModifier(code) {
			modifiers = append(modifiers, code)
		} else {
			nonModifiers = append(nonModifiers, code)
		}
	}

	if len(modifiers) < 1 {
		return KeyCombo{}, ErrInvalidKeyComboFormat
	}

	if len(nonModifiers) < 1 || len(nonModifiers) > 1 {
		return KeyCombo{}, ErrInvalidNonModifierCount
	}

	return combo, nil
}
