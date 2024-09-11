package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLogLevel(t *testing.T) {
	a := assert.New(t)
	inputs := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"DEBUG": false,
		"foo":   false,
		"bar":   false,
	}
	for input, isValid := range inputs {
		err := validateLogLevel(input)
		if isValid {
			a.NoErrorf(err, "INPUT=%s", input)
			continue
		}
		a.Errorf(err, "INPUT=%s", input)
	}
}

func TestValidateLogFormat(t *testing.T) {
	a := assert.New(t)
	inputs := map[string]bool{
		"text": true,
		"json": true,
		"TEXT": false,
		"JSON": false,
		"foo":  false,
		"bar":  false,
	}
	for input, isValid := range inputs {
		err := validateLogFormat(input)
		if isValid {
			a.NoErrorf(err, "INPUT=%s", input)
			continue
		}
		a.Errorf(err, "INPUT=%s", input)
	}
}
