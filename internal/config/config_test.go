package config

import (
	"os"
	"testing"

	"github.com/glowfi/ghkd/internal/hotkey"
	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert"
)

func loadTestConfigYAML(t *testing.T, path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Error("load test config yaml:", err)
	}
	return data
}

func TestConfigMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name         string
		cfg          Config
		cfgYAMLBytes []byte
	}{
		{
			name: "empty keybindings :POS",
			cfg: Config{
				Keybindings: []Keybinding{},
			},
			cfgYAMLBytes: loadTestConfigYAML(t, "./testdata/empty.yaml"),
		},
		{
			name: "simple run command :POS",
			cfg: Config{
				Keybindings: []Keybinding{
					{
						Name: "Open Alacritty",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTCTRL, hotkey.KEY_LEFTALT},
							Key:       hotkey.KEY_T,
							Raw:       "ctrl+alt+t",
						},
						Run: "alacritty",
					},
				},
			},
			cfgYAMLBytes: loadTestConfigYAML(t, "./testdata/run.yaml"),
		},
		{
			name: "python script :POS",
			cfg: Config{
				Keybindings: []Keybinding{
					{
						Name: "System Info",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTMETA},
							Key:       hotkey.KEY_I,
							Raw:       "super+i",
						},
						Interpreter: "python3",
						Script: `import subprocess
import os
mem = os.popen("free -h | awk '/^Mem:/ {print $3\"/\"$2}'").read().strip()
subprocess.run(["notify-send", "Memory Usage", mem])
`,
					},
				},
			},
			cfgYAMLBytes: loadTestConfigYAML(t, "./testdata/script.yaml"),
		},
		{
			name: "external file :POS",
			cfg: Config{
				Keybindings: []Keybinding{
					{
						Name: "Backup",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTMETA, hotkey.KEY_LEFTSHIFT},
							Key:       hotkey.KEY_B,
							Raw:       "super+shift+b",
						},
						File: "~/.config/hotkeysd/scripts/backup.sh",
					},
				},
			},
			cfgYAMLBytes: loadTestConfigYAML(t, "./testdata/file.yaml"),
		},
		{
			name: "multiple keybindings :POS",
			cfg: Config{
				Keybindings: []Keybinding{
					{
						Name: "Open Alacritty",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTCTRL, hotkey.KEY_LEFTALT},
							Key:       hotkey.KEY_T,
							Raw:       "ctrl+alt+t",
						},
						Run: "alacritty",
					},
					{
						Name: "System Info",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTMETA},
							Key:       hotkey.KEY_I,
							Raw:       "super+i",
						},
						Interpreter: "python3",
						Script: `import subprocess
import os
mem = os.popen("free -h | awk '/^Mem:/ {print $3\"/\"$2}'").read().strip()
subprocess.run(["notify-send", "Memory Usage", mem])
`,
					},
					{
						Name: "Backup",
						KeyCombination: hotkey.KeyCombo{
							Modifiers: []uint16{hotkey.KEY_LEFTMETA, hotkey.KEY_LEFTSHIFT},
							Key:       hotkey.KEY_B,
							Raw:       "super+shift+b",
						},
						File: "~/.config/hotkeysd/scripts/backup.sh",
					},
				},
			},
			cfgYAMLBytes: loadTestConfigYAML(t, "./testdata/multi.yaml"),
		},
	}

	for _, tt := range tests {
		// Marshal
		gotCfgYAML, gotMarshalErr := yaml.Marshal(tt.cfg)

		assert.NoError(t, gotMarshalErr, "expect no error while marshaling config")
		assert.YAMLEq(t, string(tt.cfgYAMLBytes), string(gotCfgYAML), "expect marshalled config yaml to match")

		// Unmarshal
		var gotCfg Config
		gotUnmarshalErr := yaml.Unmarshal(tt.cfgYAMLBytes, &gotCfg)

		assert.NoError(t, gotUnmarshalErr, "expect no error while unmarshaling config")
		assert.Equal(t, tt.cfg, gotCfg, "expect unmarshalled config to match")
	}
}

func TestConfig_LoadConfig(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		expectedConfig Config
		wantErr        error
	}{
		{
			name:           "should return error when no key is provided in one of the key bindings :NEG",
			configPath:     "./testdata/load_config/no_key.yaml",
			expectedConfig: Config{},
			wantErr:        hotkey.ErrInvalidKeyComboFormat,
		},
		{
			name:           "should return error when alteast one modifier key is not provided :NEG",
			configPath:     "./testdata/load_config/no_key.yaml",
			expectedConfig: Config{},
			wantErr:        hotkey.ErrInvalidKeyComboFormat,
		},
	}

	for _, tt := range tests {
		gotConfig, gotErr := LoadConfig(tt.configPath)

		assert.ErrorIs(t, gotErr, tt.wantErr, "expect error to match")
		assert.Equal(t, tt.expectedConfig, gotConfig, "expect config to match")
	}
}
