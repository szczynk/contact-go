package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"port"`
	Storage string `mapstructure:"storage"`
	Mode    string `mapstructure:"mode"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		_, ok := err.(viper.ConfigFileNotFoundError)
		if ok {
			return nil, errors.New(".env tidak ditemukan")
		}
		return nil, err
	}

	config := new(Config)
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
