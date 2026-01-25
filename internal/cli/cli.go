package cli

import (
	"flag"
	"fmt"
	"os"
)

type Command int

const (
	CommandRun Command = iota
	CommandVersion
	CommandKill
	CommandReload
	CommandBackground
)

type Options struct {
	ConfigPath string
	Command    Command
}

func Parse() (*Options, error) {
	var (
		configPath  string
		background  bool
		kill        bool
		reload      bool
		showVersion bool
	)

	// Bind both short and long flags
	flag.StringVar(&configPath, "c", "config.yaml", "config path")
	flag.StringVar(&configPath, "config", "config.yaml", "config path")

	flag.BoolVar(&background, "b", false, "background")
	flag.BoolVar(&background, "background", false, "background")

	flag.BoolVar(&kill, "k", false, "kill")
	flag.BoolVar(&kill, "kill", false, "kill")

	flag.BoolVar(&reload, "r", false, "reload")
	flag.BoolVar(&reload, "reload", false, "reload")

	flag.BoolVar(&showVersion, "v", false, "version")
	flag.BoolVar(&showVersion, "version", false, "version")

	flag.Usage = printUsage
	flag.Parse()

	opts := &Options{
		ConfigPath: configPath,
		Command:    CommandRun,
	}

	// Determine command (priority order)
	switch {
	case showVersion:
		opts.Command = CommandVersion
	case kill:
		opts.Command = CommandKill
	case reload:
		opts.Command = CommandReload
	case background:
		opts.Command = CommandBackground
	default:
		opts.Command = CommandRun
	}

	// Validate config file exists for run commands
	if opts.Command == CommandRun || opts.Command == CommandBackground {
		if err := validateConfigPath(configPath); err != nil {
			return nil, err
		}
	}

	return opts, nil
}

func validateConfigPath(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("config file not found at %s (use -c/--config to specify a valid path)", path)
		}
		return fmt.Errorf("error accessing config file %s: %w", path, err)
	}
	return nil
}

func printUsage() {
	fmt.Print(`ghkd - Go Hotkey Daemon

Usage:
  ghkd [flags]

Flags:
  -h,  --help              Prints this help message
  -c,  --config [path]     Reads the config from custom path
  -b,  --background        Runs ghkd in the background
  -k,  --kill              Gracefully kills running instances
  -r,  --reload            Reloads configuration of running instance
  -v,  --version           Prints current version
`)
}

// FilterBackgroundFlag returns args without the background flag
func FilterBackgroundFlag(args []string) []string {
	filtered := make([]string, 0, len(args))
	for _, arg := range args {
		if arg != "-b" && arg != "--background" {
			filtered = append(filtered, arg)
		}
	}
	return filtered
}
