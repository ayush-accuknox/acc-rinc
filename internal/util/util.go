package util

import (
	"fmt"
	"strings"
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
