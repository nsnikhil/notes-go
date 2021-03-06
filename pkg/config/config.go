package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	env              string
	migrationPath    string
	httpServerConfig HTTPServerConfig
	databaseConfig   DatabaseConfig
	logConfig        LogConfig
	logFileConfig    LogFileConfig
}

func (c Config) HTTPServerConfig() HTTPServerConfig {
	return c.httpServerConfig
}

func (c Config) DatabaseConfig() DatabaseConfig {
	return c.databaseConfig
}

func (c Config) LogConfig() LogConfig {
	return c.logConfig
}

func (c Config) LogFileConfig() LogFileConfig {
	return c.logFileConfig
}

func (c Config) Env() string {
	return c.env
}

func (c Config) MigrationPath() string {
	return c.migrationPath
}

func NewConfig(configFile string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(configFile)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}

	return Config{
		env:              getString("ENV"),
		migrationPath:    getString("MIGRATION_PATH"),
		httpServerConfig: newHTTPServerConfig(),
		databaseConfig:   newDatabaseConfig(),
		logConfig:        newLogConfig(),
		logFileConfig:    newLogFileConfig(),
	}
}
