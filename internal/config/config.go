package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// app
	APP_ENV  string `env:"APP_ENV"`
	APP_PORT string `env:"APP_PORT"`
	// db
	DBUSER     string `env:"DB_USER"`
	DBPASSWORD string `env:"DB_PASSWORD"`
	DBHOST     string `env:"DB_HOST"`
	DBPORT     string `env:"DB_PORT"`
	DBNAME     string `env:"DB_NAME"`
	DBSSLMODE  string `env:"DB_SSLMODE"`

	// redis
	REDIS_HOST     string `env:"REDIS_HOST"`
	REDIS_PORT     string `env:"REDIS_PORT"`
	REDIS_PASSWORD string `env:"REDIS_PASSWORD"`
	REDIS_DB       int    `env:"REDIS_DB"`
}

func Env() *Config {
	viper.SetConfigFile(".secrets/.env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}

	config := &Config{
		APP_ENV:  viper.GetString("APP_ENV"),
		APP_PORT: viper.GetString("APP_PORT"),

		DBUSER:     viper.GetString("DB_USER"),
		DBPASSWORD: viper.GetString("DB_PASSWORD"),
		DBHOST:     viper.GetString("DB_HOST"),
		DBPORT:     viper.GetString("DB_PORT"),
		DBNAME:     viper.GetString("DB_NAME"),
		DBSSLMODE:  viper.GetString("DB_SSLMODE"),

		REDIS_HOST:     viper.GetString("REDIS_HOST"),
		REDIS_PORT:     viper.GetString("REDIS_PORT"),
		REDIS_PASSWORD: viper.GetString("REDIS_PASSWORD"),
		REDIS_DB:       viper.GetInt("REDIS_DB"),
	}
	return config
}
