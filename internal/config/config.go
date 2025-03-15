package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
}

var AppConfig Config

const (
	configFileName = "config"
	configFileType = "yaml"
	configPath1    = "."
	configPath2    = ".."
	dbUserEnv      = "DB_USER"
	dbPasswordEnv  = "DB_PASSWORD"
	dbNameEnv      = "DB_NAME"
	dbHostEnv      = "DB_HOST"
	dbPortEnv      = "DB_PORT"
)

func InitConfig() {
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(configPath1) // Look for the config file in the current directory
	viper.AddConfigPath(configPath2) // Look for the config file in the parent directory
	viper.AutomaticEnv()

	// Bind environment variables
	viper.BindEnv(dbUserEnv)
	viper.BindEnv(dbPasswordEnv)
	viper.BindEnv(dbNameEnv)
	viper.BindEnv(dbHostEnv)
	viper.BindEnv(dbPortEnv)

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	AppConfig = Config{
		DBUser:     viper.GetString(dbUserEnv),
		DBPassword: viper.GetString(dbPasswordEnv),
		DBName:     viper.GetString(dbNameEnv),
		DBHost:     viper.GetString(dbHostEnv),
		DBPort:     viper.GetString(dbPortEnv),
	}
}
