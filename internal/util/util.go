package util

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/accuknox/rinc/internal/conf"
)

// IsosecLayout defines isosec timestamp.
const IsosecLayout = "20060102150405"

// GetNamespaceFromFQDN parses a Kubernetes FQDN and extracts the namespace
// from it.
func GetNamespaceFromFQDN(fqdn string) (string, error) {
	tokens := strings.Split(fqdn, ".")
	if len(tokens) < 2 {
		return "", fmt.Errorf("%q is not a fqdn", fqdn)
	}
	return tokens[1], nil
}

// NewLogger creates a new slog logger from the provided configuration.
func NewLogger(c conf.Log) *slog.Logger {
	var level slog.Level
	switch c.Level {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	opt := &slog.HandlerOptions{Level: level}
	if c.Format == "json" {
		return slog.New(slog.NewJSONHandler(os.Stderr, opt))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, opt))
}
