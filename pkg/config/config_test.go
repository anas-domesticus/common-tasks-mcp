package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConfig struct {
	ListenPort    int      `env:"LISTEN_PORT" yaml:"listenPort" default:"8080"`
	DefaultField  string   `env:"DEFAULT_FIELD" yaml:"defaultField" default:"foo"`
	RequiredField string   `env:"REQUIRED_FIELD" yaml:"requiredField" required:"true"`
	Tags          []string `env:"TAGS" yaml:"tags"`
}

// testConfigWithNamedInline tests named inline structs
type testConfigWithNamedInline struct {
	ListenPort int            `env:"LISTEN_PORT" yaml:"listenPort" default:"8080"`
	Database   DatabaseConfig `yaml:"database,inline"`
	Cache      CacheConfig    `yaml:"cache,inline"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST" yaml:"host" default:"localhost"`
	Port     int    `env:"DB_PORT" yaml:"port" default:"5432"`
	Password string `env:"DB_PASSWORD" yaml:"password" required:"true"`
}

type CacheConfig struct {
	Address  string `env:"CACHE_ADDRESS" yaml:"address" default:"redis:6379"`
	Password string `env:"CACHE_PASSWORD" yaml:"password" default:""`
	DB       int    `env:"CACHE_DB" yaml:"db" default:"0"`
}

func TestGetConfigFromEnvVars(t *testing.T) {
	testCases := []struct {
		name    string
		envVars map[string]string
		want    testConfig
		wantErr bool
	}{
		{
			name: "All defaults, except for required field",
			envVars: map[string]string{
				"REQUIRED_FIELD": "bar",
			},
			want: testConfig{
				ListenPort:    8080,
				DefaultField:  "foo",
				RequiredField: "bar",
			},
			wantErr: false,
		},
		{
			name: "ListenPort from env",
			envVars: map[string]string{
				"LISTEN_PORT":    "9000",
				"REQUIRED_FIELD": "bar",
			},
			want: testConfig{
				ListenPort:    9000,
				DefaultField:  "foo",
				RequiredField: "bar",
			},
			wantErr: false,
		},
		{
			name: "Override default field from env",
			envVars: map[string]string{
				"DEFAULT_FIELD":  "custom",
				"REQUIRED_FIELD": "bar",
			},
			want: testConfig{
				ListenPort:    8080,
				DefaultField:  "custom",
				RequiredField: "bar",
			},
			wantErr: false,
		},
		{
			name: "Invalid ListenPort",
			envVars: map[string]string{
				"LISTEN_PORT":    "invalid",
				"REQUIRED_FIELD": "bar",
			},
			want:    testConfig{},
			wantErr: true,
		},
		{
			name:    "Missing required field",
			envVars: map[string]string{},
			want:    testConfig{},
			wantErr: true,
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

			var config testConfig
			err := GetConfigFromEnvVars(&config)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, config)
		})
	}
}

func TestGetConfig(t *testing.T) {
	yamlData := `
listenPort: 9000
requiredField: bar
`
	testCases := []struct {
		name            string
		filepath        string
		allowFileErrors bool
		envVars         map[string]string
		fileContent     string
		want            testConfig
		wantErr         bool
	}{
		{
			name:            "YAML file overridden by env vars",
			filepath:        "config.yaml",
			allowFileErrors: false,
			envVars: map[string]string{
				"LISTEN_PORT": "8001",
			},
			fileContent: yamlData,
			want: testConfig{
				ListenPort:    8001,
				DefaultField:  "foo",
				RequiredField: "bar",
			},
			wantErr: false,
		},
		{
			name:            "Env vars override defaults with missing file",
			filepath:        "missing.yaml",
			allowFileErrors: true,
			envVars: map[string]string{
				"LISTEN_PORT":    "8001",
				"DEFAULT_FIELD":  "custom",
				"REQUIRED_FIELD": "bar",
			},
			want: testConfig{
				ListenPort:    8001,
				DefaultField:  "custom",
				RequiredField: "bar",
			},
			wantErr: false,
		},
		{
			name:            "File error without allowFileErrors",
			filepath:        "missing.yaml",
			allowFileErrors: false,
			want:            testConfig{},
			wantErr:         true,
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

			if tc.fileContent != "" {
				if err := os.WriteFile(tc.filepath, []byte(tc.fileContent), 0644); err != nil {
					t.Fatalf("failed to write temp file: %v", err)
				}
				defer os.Remove(tc.filepath)
			}

			var config testConfig
			err := GetConfig(&config, tc.filepath, tc.allowFileErrors)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, config)
		})
	}
}

func TestGetConfigFromEnvVars_NamedInlineStruct(t *testing.T) {
	testCases := []struct {
		name    string
		envVars map[string]string
		want    testConfigWithNamedInline
		wantErr bool
	}{
		{
			name: "Named inline struct with defaults",
			envVars: map[string]string{
				"DB_PASSWORD": "secret123",
			},
			want: testConfigWithNamedInline{
				ListenPort: 8080,
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					Password: "secret123",
				},
				Cache: CacheConfig{
					Address:  "redis:6379",
					Password: "",
					DB:       0,
				},
			},
			wantErr: false,
		},
		{
			name: "Named inline struct with env overrides",
			envVars: map[string]string{
				"DB_HOST":        "postgres.example.com",
				"DB_PORT":        "5433",
				"DB_PASSWORD":    "secret123",
				"CACHE_ADDRESS":  "redis.example.com:6380",
				"CACHE_PASSWORD": "cache_secret",
				"CACHE_DB":       "1",
			},
			want: testConfigWithNamedInline{
				ListenPort: 8080,
				Database: DatabaseConfig{
					Host:     "postgres.example.com",
					Port:     5433,
					Password: "secret123",
				},
				Cache: CacheConfig{
					Address:  "redis.example.com:6380",
					Password: "cache_secret",
					DB:       1,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing required field in named inline struct",
			envVars: map[string]string{
				"DB_HOST": "localhost",
			},
			want:    testConfigWithNamedInline{},
			wantErr: true,
		},
		{
			name: "Invalid type conversion in named inline struct",
			envVars: map[string]string{
				"DB_PASSWORD": "secret123",
				"DB_PORT":     "invalid_port",
			},
			want:    testConfigWithNamedInline{},
			wantErr: true,
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

			var config testConfigWithNamedInline
			err := GetConfigFromEnvVars(&config)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.want, config)
		})
	}
}
