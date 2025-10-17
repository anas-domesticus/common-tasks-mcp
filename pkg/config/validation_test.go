package config

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type validationTestConfig struct {
	ListenPort    int      `env:"LISTEN_PORT" yaml:"listenPort" default:"8080"`
	RequiredField string   `env:"REQUIRED_FIELD" yaml:"requiredField" required:"true"`
	Tags          []string `env:"TAGS" yaml:"tags"`
}

// Validate implements custom validation for testing
func (c validationTestConfig) Validate() error {
	var errs []error
	if c.ListenPort < 1 || c.ListenPort > 65535 {
		errs = append(errs, fmt.Errorf("listenPort must be between 1-65535, got %d", c.ListenPort))
	}
	return errors.Join(errs...)
}

func TestValidation(t *testing.T) {
	testCases := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid port",
			envVars: map[string]string{
				"LISTEN_PORT":    "8000",
				"REQUIRED_FIELD": "bar",
			},
			wantErr: false,
		},
		{
			name: "Valid custom port",
			envVars: map[string]string{
				"LISTEN_PORT":    "9000",
				"REQUIRED_FIELD": "bar",
			},
			wantErr: false,
		},
		{
			name: "Invalid port - too low",
			envVars: map[string]string{
				"LISTEN_PORT":    "0",
				"REQUIRED_FIELD": "bar",
			},
			wantErr: true,
			errMsg:  "listenPort must be between 1-65535",
		},
		{
			name: "Invalid port - too high",
			envVars: map[string]string{
				"LISTEN_PORT":    "70000",
				"REQUIRED_FIELD": "bar",
			},
			wantErr: true,
			errMsg:  "listenPort must be between 1-65535",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for key, value := range tc.envVars {
				os.Setenv(key, value)
			}
			defer func() {
				for key := range tc.envVars {
					os.Unsetenv(key)
				}
			}()

			var config validationTestConfig
			err := GetConfigFromEnvVars(&config)

			// Debug output
			t.Logf("Config: %+v", config)
			t.Logf("Error: %v", err)

			if tc.wantErr {
				assert.Error(t, err)
				if err != nil {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSliceSupport(t *testing.T) {
	os.Setenv("TAGS", "tag1,tag2, tag3")
	os.Setenv("REQUIRED_FIELD", "bar")
	defer func() {
		os.Unsetenv("TAGS")
		os.Unsetenv("REQUIRED_FIELD")
	}()

	var config validationTestConfig
	err := GetConfigFromEnvVars(&config)
	assert.NoError(t, err)
	assert.Equal(t, []string{"tag1", "tag2", "tag3"}, config.Tags)
}
