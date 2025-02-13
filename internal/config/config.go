package config

import (
	"rashpor.com/todolist/internal/config/database"
)

type Config struct {
	Database *database.Config
}

func NewDBConfig() *Config {
	return &Config{
		Database: &database.Config{
			User:     getEnv("DB_CONFIG_USER", "skillsrock"),
			Password: getEnv("DB_CONFIG_PASSWORD", "ForSkillsRock"),
			Host:     getEnv("DB_CONFIG_HOST", "localhost"),
			Port:     getEnvAsInt("DB_CONFIG_PORT", 5432),
			DbName:   getEnv("DB_CONFIG_DBNAME", "todolist"),
		},
	}
}
