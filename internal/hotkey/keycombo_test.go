package hotkey

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyCombo_ParseKeyCombo(t *testing.T) {
	tests := []struct {
		name             string
		inputKeyCombo    string
		expectedKeyCombo KeyCombo
		wantErr          error
	}{
		{
			name:             "should return error when no key is provided :NEG",
			inputKeyCombo:    "",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrInvalidKeyComboFormat,
		},
		{
			name:             "should return error when combination of modifier and non-modifier is not provided :NEG",
			inputKeyCombo:    "ctrl",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrInvalidKeyComboFormat,
		},
		{
			name:             "should return error when invalid modifier key is provided :NEG",
			inputKeyCombo:    "ctrl+altfoo+a",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrUnknownKey,
		},
		{
			name:             "should return error when invalid non-modifier key is provided :NEG",
			inputKeyCombo:    "ctrl+alt+foo",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrUnknownKey,
		},
		{
			name:             "should return error when more than one non-modifier keys are provided :NEG",
			inputKeyCombo:    "ctrl+alt+b+c",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrInvalidNonModifierCount,
		},
		{
			name:             "should return error when more than no modifier keys are provided :NEG",
			inputKeyCombo:    "a+b+c",
			expectedKeyCombo: KeyCombo{},
			wantErr:          ErrInvalidKeyComboFormat,
		},
		{
			name:          "should parse key combo successfully :POS",
			inputKeyCombo: "ctrl+SHIFT+b",
			expectedKeyCombo: KeyCombo{
				Modifiers: []uint16{
					KEY_LEFTCTRL,
					KEY_LEFTSHIFT,
				},
				Key: KEY_B,
				Raw: "ctrl+SHIFT+b",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		gotKeyCombo, gotErr := ParseKeyCombo(tt.inputKeyCombo)

		if tt.wantErr != nil {
			assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
			return
		}

		assert.NoError(t, gotErr, "expect no error while parsing key combination")
		assert.Equal(t, tt.expectedKeyCombo, gotKeyCombo, "expect key combination to match")
	}
}

func TestKeyCombo_String(t *testing.T) {
	tests := []struct {
		name             string
		inputKeyCombo    KeyCombo
		expectedKeyCombo string
	}{
		{
			name: "should return key combo invalid string representation for invalid key combo input :NEG",
			inputKeyCombo: KeyCombo{
				Modifiers: []uint16{
					10_000, 20_000,
				},
				Key: 30_000,
			},
			expectedKeyCombo: "",
		},
		{
			name: "should return key combo string representation for valid key combo input :POS",
			inputKeyCombo: KeyCombo{
				Modifiers: []uint16{
					KEY_LEFTCTRL, KEY_LEFTALT,
				},
				Key: KEY_S,
			},
			expectedKeyCombo: "ctrl+alt+s",
		},
	}

	for _, tt := range tests {
		gotKeyCombo := tt.inputKeyCombo.String()

		assert.Equal(t, tt.expectedKeyCombo, gotKeyCombo, "expect key combo string representation to match")
	}
}
