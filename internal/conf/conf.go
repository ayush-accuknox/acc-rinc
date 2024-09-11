package conf

import (
	"fmt"
	"os"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	flag "github.com/spf13/pflag"
)

var defaultConfig = "/etc/reporter/config.yaml"

// C contains all configuration data that can be passed to the reporter.
type C struct {
	// Log contains configuration for logs.
	Log Log `koanf:"log"`
	// Output is the path to the reports output directory.
	Output string `koanf:"output"`
	// RabbitMQ contains the rabbitmq configuration.
	RabbitMQ RabbitMQ `koanf:"rabbitmq"`
}

// New creates a configuration using the provided arguments and config file.
func New(args ...string) (*C, error) {
	k := koanf.New(".")

	err := k.Load(confmap.Provider(map[string]any{
		"log.level":                    "info",
		"log.format":                   "text",
		"output":                       "reports",
		"rabbitmq.enable":              false,
		"rabbitmq.headlessSvcAddr":     "",
		"rabbitmq.management.url":      "",
		"rabbitmq.management.username": "",
		"rabbitmq.management.password": "",
	}, "."), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load default configuration: %w", err)
	}

	f := parseFlags(args)
	confF, err := f.GetString("conf")
	if err != nil {
		return nil, fmt.Errorf("failed to parse flags: %w", err)
	}

	err = k.Load(file.Provider(confF), yaml.Parser())
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	conf := new(C)
	err = k.Unmarshal("", conf)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal configuration: %w", err)
	}

	return conf, nil
}

func parseFlags(args []string) *flag.FlagSet {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Print(f.FlagUsages())
		os.Exit(0)
	}
	f.String("conf", defaultConfig, "path to config file")
	f.Parse(args)
	return f
}
