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
	}

	for _, tt := range tests {
		gotCfgYAML, gotMarshalErr := yaml.Marshal(tt.cfg)

		assert.ErrorIs(t, gotMarshalErr, tt.wantMarshalErr, "expect error to match while marshaling config")
		assert.Equal(t, string(tt.cfgYAMLBytes), string(gotCfgYAML), "expect marshalled config yaml to match")

		var gotCfg Config
		gotUnmarshalErr := yaml.Unmarshal(tt.cfgYAMLBytes, &gotCfg)

		assert.ErrorIs(t, gotUnmarshalErr, tt.wantUnMarshalErr, "expect error to match while unmarshaling config")
		assert.Equal(t, tt.cfg, gotCfg, "expect unmarshalled config to match")
	}
}
