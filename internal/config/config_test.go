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
		name             string
		cfg              Config
		cfgYAMLBytes     []byte
		wantMarshalErr   error
		wantUnMarshalErr error
	}{
		{
			name: "empty keybindings :POS",
			cfg: Config{
				Keybindings: []Keybinding{},
			},
			cfgYAMLBytes:     loadTestConfigYAML(t, "./testdata/empty.yaml"),
			wantMarshalErr:   nil,
			wantUnMarshalErr: nil,
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
			cfgYAMLBytes:     loadTestConfigYAML(t, "./testdata/run.yaml"),
			wantMarshalErr:   nil,
			wantUnMarshalErr: nil,
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
			cfgYAMLBytes:     loadTestConfigYAML(t, "./testdata/script.yaml"),
			wantMarshalErr:   nil,
			wantUnMarshalErr: nil,
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
			cfgYAMLBytes:     loadTestConfigYAML(t, "./testdata/file.yaml"),
			wantMarshalErr:   nil,
			wantUnMarshalErr: nil,
		},
	}

	for _, tt := range tests {
		// Marshal
		gotCfgYAML, gotMarshalErr := yaml.Marshal(tt.cfg)

		assert.ErrorIs(t, gotMarshalErr, tt.wantMarshalErr, "expect error to match while marshaling config")
		assert.YAMLEq(t, string(tt.cfgYAMLBytes), string(gotCfgYAML), "expect marshalled config yaml to match semantically")

		// Unmarshal
		var gotCfg Config
		gotUnmarshalErr := yaml.Unmarshal(tt.cfgYAMLBytes, &gotCfg)

		assert.ErrorIs(t, gotUnmarshalErr, tt.wantUnMarshalErr, "expect error to match while unmarshaling config")
		assert.Equal(t, tt.cfg, gotCfg, "expect unmarshalled config to match")
	}
}
