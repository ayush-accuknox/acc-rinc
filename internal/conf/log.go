package conf

// Log contains configuration for logs.
type Log struct {
	// Level is the log level.
	// Possible values: "debug", "info", "warn", "error".
	//
	// Default: "info"
	Level string `koanf:"level"`
	// Format specifies the format of the logs.
	// Possible values: "text", "json"
	//
	// Default: "text"
	Format string `koanf:"format"`
}
