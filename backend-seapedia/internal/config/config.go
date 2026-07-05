package config

import (
	"backend-seapedia/internal/model"
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() *model.Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Warning: file .env tidak ditemukan")
	}

	return &model.Config{
		AppPort:    viper.GetString("APP_PORT"),
		DBHost:     viper.GetString("DB_HOST"),
		DBPort:     viper.GetString("DB_PORT"),
		DBUser:     viper.GetString("DB_USER"),
		DBPassword: viper.GetString("DB_PASSWORD"),
		DBName:     viper.GetString("DB_NAME"),
		JWTSecret:  viper.GetString("JWT_SECRET"),
	}
}
