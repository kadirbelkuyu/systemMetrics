package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	fmt.Println("Searching for config in the following paths:")
	for _, path := range viper.GetStringSlice("configPaths") {
		fmt.Println(path)
	}

	var cfg Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file: %v\n", err)
		return nil, err
	}
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
