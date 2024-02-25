package config

import (
	"fmt"

	"github.com/platatest/internal/config"
	"github.com/spf13/viper"
)

func Parse() (config.Config, error) {
	var c config.Config

	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		return c, fmt.Errorf("%w : %w", ReadErr, err)
	}

	err = v.Unmarshal(&c)
	if err != nil {
		return c, fmt.Errorf("%w : %w", UnmarshalErr, err)
	}
	return c, err
}
