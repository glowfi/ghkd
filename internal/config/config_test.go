package config

import (
	"os"
	"testing"

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
						Keys: "ctrl+alt+t",
						Run:  "alacritty",
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
						Name:        "System Info",
						Keys:        "super+i",
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
						Keys: "super+shift+b",
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
						Keys: "ctrl+alt+t",
						Run:  "alacritty",
					},
					{
						Name:        "System Info",
						Keys:        "super+i",
						Interpreter: "python3",
						Script: `import subprocess
import os
mem = os.popen("free -h | awk '/^Mem:/ {print $3\"/\"$2}'").read().strip()
subprocess.run(["notify-send", "Memory Usage", mem])
`,
					},
					{
						Name: "Backup",
						Keys: "super+shift+b",
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

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		inputConfig Config
		wantErr     bool
	}{
		{
			name: "should return error when no name is provided :NEG",
			inputConfig: Config{
				Keybindings: []Keybinding{
					{
						Name: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when no key is provided :NEG",
			inputConfig: Config{
				Keybindings: []Keybinding{
					{
						Name: "foo",
						Keys: "",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when modifier key is not provided :NEG",
			inputConfig: Config{
				Keybindings: []Keybinding{
					{
						Name: "foo",
						Keys: "a",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "should return error when non-modifier key is not provided :NEG",
			inputConfig: Config{
				Keybindings: []Keybinding{
					{
						Name: "foo",
						Keys: "ctrl",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		gotErr := tt.inputConfig.Validate()

		if tt.wantErr {
			assert.Error(t, gotErr, "expect error while validating config")
			return
		}

		assert.NoError(t, gotErr, "expect no error while validating config")
	}
}
