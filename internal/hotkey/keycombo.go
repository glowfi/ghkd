package hotkey

import (
	"errors"
	"fmt"
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

		if IsModifier(part) {
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

	combo.Modifiers = modifiers
	combo.Key = nonModifiers[0]

	return combo, nil
}

func (kc KeyCombo) String() string {
	if kc.Raw != "" {
		return kc.Raw
	}

	var parts []string
	for _, mod := range kc.Modifiers {
		name, found := LookupKeyName(mod)
		if found {
			parts = append(parts, name)
		}
	}
	if kc.Key != 0 {
		name, found := LookupKeyName(kc.Key)
		if found {
			parts = append(parts, name)
		}
	}
	return strings.Join(parts, "+")
}

// UnmarshalYAML implements custom YAML unmarshaling
func (kc *KeyCombo) UnmarshalYAML(unmarshal func(any) error) error {
	var raw string
	if err := unmarshal(&raw); err != nil {
		return err
	}

	parsed, err := ParseKeyCombo(raw)
	if err != nil {
		return fmt.Errorf("parse key combination '%s': %w", raw, err)
	}

	*kc = parsed
	return nil
}

// MarshalYAML implements custom YAML marshaling
func (kc KeyCombo) MarshalYAML() (any, error) {
	return kc.String(), nil
}

// Matches checks if pressed keys match this combo
func (kc KeyCombo) Matches(pressed []uint16) bool {
	n := len(kc.Modifiers) + 1
	if len(pressed) != n {
		return false
	}

	if kc.Key != pressed[n-1] {
		return false
	}

	for idx, val := range kc.Modifiers {
		if val != pressed[idx] {
			return false
		}
	}

	return true
}
