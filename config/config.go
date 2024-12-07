// Package config contains all the config
package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

const (
	defaultConfigFolder  = "./deployment/development"
	modifiedConfigFolder = "./development"
)

func New() (Schema, error) {
	log.Info("Loading config")

	var cfg Schema

	viper.AddConfigPath(defaultConfigFolder)
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("error reading config: %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error unmarshalling config: %w", err)
	}

	validate := validator.New()
	err = validate.Struct(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("error validating config: %w", err)
	}

	return cfg, nil
}
