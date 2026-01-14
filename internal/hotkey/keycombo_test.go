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
