package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"
)

type Config struct {
	Port         int    `env:"PORT"`
	JWT_SECRET   string `env:"JWT_SECRET"`
	POSTGRES_DSN string `env:"POSTGRES_DSN"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	cfgVal := reflect.ValueOf(cfg).Elem()
	cfgType := cfgVal.Type()

	durationType := reflect.TypeOf(time.Duration(0))

	for i := 0; i < cfgVal.NumField(); i++ {
		field := cfgVal.Field(i)
		fieldType := cfgType.Field(i)
		tag := fieldType.Tag.Get("env")

		if tag == "" {
			continue
		}

		envVal := os.Getenv(tag)
		if envVal == "" {
			return cfg, fmt.Errorf("environment variable %s is not set", tag)
		}

		kind := field.Kind()

		switch {
		case reflect.String == kind:
			field.SetString(envVal)

		case field.Type() == durationType:
			secVal, err := strconv.ParseInt(envVal, 10, 64)
			if err != nil {
				durVal, err2 := time.ParseDuration(envVal)
				if err2 != nil {
					return cfg, fmt.Errorf("environment variable %s is not a valid duration: %v", tag, err)
				}
				field.SetInt(int64(durVal))
			} else {
				field.SetInt(int64(time.Duration(secVal) * time.Second))
			}

		case reflect.Int == kind:
			intVal, err := strconv.Atoi(envVal)
			if err != nil {
				return cfg, fmt.Errorf("environment variable %s is not a valid integer", tag)
			}
			field.SetInt(int64(intVal))

		default:
			return cfg, fmt.Errorf("unsupported type %s for field %s", field.Type(), fieldType.Name)
		}
	}
	return cfg, nil
}
