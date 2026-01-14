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
		wantErr          bool
	}{
		{
			name:             "should return error when no key is provided :NEG",
			inputKeyCombo:    "",
			expectedKeyCombo: KeyCombo{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		gotKeyCombo, gotErr := ParseKeyCombo(tt.inputKeyCombo)

		if tt.wantErr {
			assert.Error(t, gotErr, "expect error while parsing key combination")
			return
		}

		assert.NoError(t, gotErr, "expect no error while parsing key combination")
		assert.Equal(t, tt.expectedKeyCombo, gotKeyCombo, "expect key combination to match")
	}
}
