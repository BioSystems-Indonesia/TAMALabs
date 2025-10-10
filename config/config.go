// Package config contains all the config
package config

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/constant"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func New(db *gorm.DB) (Schema, error) {
	slog.Info("loading config")

	var configs []entity.Config
	err := db.Find(&configs).Error
	if err != nil {
		return Schema{}, err
	}

	v := viper.New()

	var mapping = map[string]string{}
	for _, config := range configs {
		mapping[config.ID] = config.Value
		v.Set(config.ID, config.Value)
	}

	cfg := Schema{}
	err = v.Unmarshal(&cfg)
	if err != nil {
		return Schema{}, fmt.Errorf("error unmarshalling config: %w", err)
	}

	InjectRuntime(&cfg)

	validate := validator.New()
	err = validate.Struct(&cfg)
	if err != nil {
		return Schema{}, fmt.Errorf("error validating config: %w", err)
	}

	return cfg, nil
}

func InjectRuntime(cfg *Schema) {
	cfg.Version = constant.AppVersion
	cfg.Revision = time.Now().Format("20060102_150405")
}
