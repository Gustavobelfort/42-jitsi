package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// RequiredConfigError is an error type that report required configuration fields that are missing a value.
type RequiredConfigError struct {
	missing []string
}

// Error formats the RequiredConfigError dynamically depending on the missing fields list.
func (err *RequiredConfigError) Error() string {
	return fmt.Sprintf("missing required configuration: %s", strings.Join(err.missing, ","))
}

// checkRequired takes as parameters the required fields. If a field's value is missing, the field is added
// to the missing fields list and an error is returned.
//
// Otherwise, nil is returned.
func checkRequired(required ...string) error {
	missing := make([]string, 0)

	for _, req := range required {
		if !viper.IsSet(req) {
			missing = append(missing, req)
		}
	}
	if len(missing) == 0 {
		return nil
	}
	return &RequiredConfigError{missing: missing}
}
