# Config

Type-safe configuration loading from YAML files and environment variables with validation.

## Purpose
Generic configuration loader supporting struct tags for environment variables, defaults, required fields, and custom validation with proper precedence handling.

## Features
- Load configuration from YAML files
- Override with environment variables
- Support for default values
- Required field validation
- Custom validation via `Validator` interface
- Support for embedded structs and inline structs
- Type conversion for strings, ints, floats, bools, and string slices

## Example
```go
package main

import "common-tasks-mcp/pkg/config"

type MyConfig struct {
    ListenPort  int    `env:"LISTEN_PORT" yaml:"listenPort" default:"8080"`
    DatabaseURL string `env:"DATABASE_URL" yaml:"databaseUrl" required:"true"`
    MaxRetries  int    `env:"MAX_RETRIES" yaml:"maxRetries" default:"3"`
}

func main() {
    var cfg MyConfig
    err := config.GetConfig(&cfg, "config.yaml", true)
    if err != nil {
        panic(err)
    }
    // Use cfg...
}
```

## Struct Tags
- `env:"ENV_VAR_NAME"` - Environment variable name
- `yaml:"fieldName"` - YAML field name
- `required:"true"` - Field is required
- `default:"value"` - Default value if not set

## Custom Validation
Implement the `Validator` interface to add custom validation:

```go
func (c MyConfig) Validate() error {
    if c.ListenPort < 1 || c.ListenPort > 65535 {
        return fmt.Errorf("listenPort must be between 1-65535, got %d", c.ListenPort)
    }
    return nil
}
```
