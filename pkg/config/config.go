package config

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Validator interface allows config structs to implement custom validation logic.
// If a config struct implements this interface, validation will be automatically
// called after loading configuration from files and environment variables.
type Validator interface {
	Validate() error
}

func processFields(val reflect.Value, typeOfT reflect.Type) (map[string]bool, error) {
	setFields := make(map[string]bool)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typeOfT.Field(i)

		// Handle embedded structs (both anonymous and named with inline tag)
		if field.Kind() == reflect.Struct && (fieldType.Anonymous || strings.Contains(fieldType.Tag.Get("yaml"), "inline")) {
			embeddedSetFields, err := processFields(field, fieldType.Type)
			if err != nil {
				return nil, err
			}
			// Merge embedded set fields
			for k, v := range embeddedSetFields {
				setFields[k] = v
			}
		} else {
			tag := fieldType.Tag.Get("env")
			if tag != "" {
				envVal := os.Getenv(tag)
				if envVal == "" {
					continue
				}

				// Mark this field as set from environment
				setFields[fieldType.Name] = true

				// Set the value to the field based on its type
				switch field.Kind() {
				case reflect.String:
					field.SetString(envVal)
				case reflect.Int:
					intVal, err := strconv.Atoi(envVal)
					if err != nil {
						return nil, fmt.Errorf("failed to convert %s to int: %v", envVal, err)
					}
					field.SetInt(int64(intVal))
				case reflect.Float64:
					floatVal, err := strconv.ParseFloat(envVal, 64)
					if err != nil {
						return nil, fmt.Errorf("failed to convert %s to float64: %v", envVal, err)
					}
					field.SetFloat(floatVal)
				case reflect.Float32:
					floatVal, err := strconv.ParseFloat(envVal, 32)
					if err != nil {
						return nil, fmt.Errorf("failed to convert %s to float32: %v", envVal, err)
					}
					field.SetFloat(floatVal)
				case reflect.Bool:
					boolVal, err := strconv.ParseBool(envVal)
					if err != nil {
						return nil, fmt.Errorf("failed to convert %s to bool: %v", envVal, err)
					}
					field.SetBool(boolVal)
				case reflect.Slice:
					// Handle string slices (comma-separated values)
					if field.Type().Elem().Kind() == reflect.String {
						values := strings.Split(envVal, ",")
						slice := reflect.MakeSlice(field.Type(), len(values), len(values))
						for i, v := range values {
							slice.Index(i).SetString(strings.TrimSpace(v))
						}
						field.Set(slice)
					} else {
						return nil, fmt.Errorf("unsupported slice type %s", field.Type())
					}
				default:
					return nil, fmt.Errorf("unsupported kind %s", field.Kind())
				}
			}
		}
	}
	return setFields, nil
}

func checkRequiredAndDefaults(val reflect.Value, typeOfT reflect.Type, setFields map[string]bool) error {
	var errs []error
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typeOfT.Field(i)

		// Handle embedded structs (both anonymous and named with inline tag)
		if field.Kind() == reflect.Struct && (fieldType.Anonymous || strings.Contains(fieldType.Tag.Get("yaml"), "inline")) {
			if err := checkRequiredAndDefaults(field, fieldType.Type, setFields); err != nil {
				errs = append(errs, err)
			}
		} else {
			fieldRequired := false
			requiredTag := fieldType.Tag.Get("required")
			if strings.ToLower(requiredTag) == "true" || strings.ToLower(requiredTag) == "1" {
				fieldRequired = true
			}
			defaultTag := fieldType.Tag.Get("default")
			if fieldRequired && defaultTag != "" { // ignoring required tag if default is set
				fieldRequired = false
			}

			if field.IsZero() && fieldRequired {
				envTag := fieldType.Tag.Get("env")
				yamlTag := fieldType.Tag.Get("yaml")
				errs = append(errs, fmt.Errorf("required field env:%s / yaml:%s is missing", envTag, yamlTag))
				continue
			}

			// Only apply defaults if the field wasn't explicitly set from environment
			if field.IsZero() && defaultTag != "" && !setFields[fieldType.Name] {
				switch field.Kind() {
				case reflect.String:
					field.SetString(defaultTag)
				case reflect.Int:
					intVal, err := strconv.Atoi(defaultTag)
					if err != nil {
						errs = append(errs, fmt.Errorf("failed to convert %s to int: %v", defaultTag, err))
					}
					field.SetInt(int64(intVal))
				case reflect.Float64:
					floatVal, err := strconv.ParseFloat(defaultTag, 64)
					if err != nil {
						errs = append(errs, fmt.Errorf("failed to convert %s to float64: %v", defaultTag, err))
					}
					field.SetFloat(floatVal)
				case reflect.Float32:
					floatVal, err := strconv.ParseFloat(defaultTag, 32)
					if err != nil {
						errs = append(errs, fmt.Errorf("failed to convert %s to float32: %v", defaultTag, err))
					}
					field.SetFloat(floatVal)
				case reflect.Bool:
					boolVal, err := strconv.ParseBool(defaultTag)
					if err != nil {
						errs = append(errs, fmt.Errorf("failed to convert %s to bool: %v", defaultTag, err))
					}
					field.SetBool(boolVal)
				default:
					errs = append(errs, fmt.Errorf("unsupported kind %s", field.Kind()))
				}
			}
		}
	}
	return errors.Join(errs...)
}

// GetConfigFromEnvVars loads configuration from environment variables only.
// It processes struct tags: env, default, required.
// Example usage:
//
//	var cfg MyConfig
//	err := GetConfigFromEnvVars(&cfg)
func GetConfigFromEnvVars[T any](dest *T) error {
	val := reflect.ValueOf(dest).Elem()
	typeOfT := val.Type()
	setFields, err := processFields(val, typeOfT)
	if err != nil {
		return err
	}
	err = checkRequiredAndDefaults(val, typeOfT, setFields)
	if err != nil {
		*dest = reflect.New(reflect.TypeOf(dest).Elem()).Elem().Interface().(T) // resets config to empty
		return err
	}

	// Run custom validation if the type implements Validator
	if validator, ok := any(*dest).(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}

// GetConfig loads configuration from YAML file first, then overlays environment variables.
// If filepath is empty, only environment variables are used.
// If allowFileErrors is true, file read/parse errors fallback to env vars only.
// Example usage:
//
//	var cfg MyConfig
//	err := GetConfig(&cfg, "config.yaml", true)
func GetConfig[T any](dest *T, filepath string, allowFileErrors bool) error {
	if filepath == "" {
		return GetConfigFromEnvVars(dest)
	}
	data, err := os.ReadFile(filepath)
	if err != nil {
		if allowFileErrors {
			return GetConfigFromEnvVars(dest)
		}
		return fmt.Errorf("failed to read file: %w", err)
	}
	if err := yaml.Unmarshal(data, dest); err != nil {
		if allowFileErrors {
			return GetConfigFromEnvVars(dest)
		}
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	err = GetConfigFromEnvVars(dest)
	if err != nil {
		return err
	}

	// Run custom validation if the type implements Validator
	if validator, ok := any(*dest).(Validator); ok {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}
	}

	return nil
}
