package app

const Version = "1.0.1"

type Config struct {
	InputDir    string
	CfgPath     string
	PidFilePath string
}

func NewConfig(InputDir, configPath, PidFilePath string) *Config {
	return &Config{
		InputDir:    InputDir,
		CfgPath:     configPath,
		PidFilePath: PidFilePath,
	}
}
