package util

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	User                            string `mapstructure:"user"`
	Host                            string `mapstructure:"host"`
	Port                            string `mapstructure:"port"`
	DBName                          string `mapstructure:"dbname"`
	Password                        string `mapstructure:"password"`
	SSLMode                         string `mapstructure:"sslMode"`
	ConnectionTimeout               int    `mapstructure:"connectionTimeout"`
	StatementTimeout                int    `mapstructure:"statementTimeout"`
	IdleInTransactionSessionTimeout int    `mapstructure:"idleInTransactionSessionTimeout"`
}

// NewDBStringFromConfig build database connection string from config file.
func NewDBStringFromConfig(config *viper.Viper) (string, error) {
	var allConfig struct {
		Database DBConfig `mapstructure:"database"`
	}
	if err := config.Unmarshal(&allConfig); err != nil {
		return "", fmt.Errorf("cannot unmarshal db config: %w", err)
	}

	return NewDBStringFromDBConfig(allConfig.Database)
}

func NewDBStringFromDBConfig(config DBConfig) (string, error) {
	var dbParams []string
	dbParams = append(dbParams, fmt.Sprintf("user=%s", config.User))
	dbParams = append(dbParams, fmt.Sprintf("host=%s", config.Host))
	dbParams = append(dbParams, fmt.Sprintf("port=%s", config.Port))
	dbParams = append(dbParams, fmt.Sprintf("dbname=%s", config.DBName))
	if password := config.Password; password != "" {
		dbParams = append(dbParams, fmt.Sprintf("password=%s", password))
	}
	dbParams = append(dbParams, fmt.Sprintf("sslmode=%s",
		config.SSLMode))
	dbParams = append(dbParams, fmt.Sprintf("connect_timeout=%d",
		config.ConnectionTimeout))
	dbParams = append(dbParams, fmt.Sprintf("statement_timeout=%d",
		config.StatementTimeout))
	dbParams = append(dbParams, fmt.Sprintf("idle_in_transaction_session_timeout=%d",
		config.IdleInTransactionSessionTimeout))

	return strings.Join(dbParams, " "), nil
}
