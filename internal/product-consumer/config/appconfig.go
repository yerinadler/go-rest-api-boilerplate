package config

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (*AppConfig, error) {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./internal/product-consumer")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var appConfig AppConfig
	if err := viper.Unmarshal(&appConfig); err != nil {
		return nil, err
	}

	return &appConfig, nil
}
