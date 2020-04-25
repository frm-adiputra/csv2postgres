package pipeline

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// DBConfig specifies database connection configurations.
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     int32
	Db       string
	SSLMode  string `yaml:"sslmode"`
}

// NewDBConfig creates a new database configuration from YAML file.
func NewDBConfig(configFile string) (DBConfig, error) {
	f, err := os.Open(configFile)
	if err != nil {
		return DBConfig{}, err
	}
	defer f.Close()

	t := DBConfig{}
	d := yaml.NewDecoder(f)
	err = d.Decode(&t)
	if err != nil {
		return DBConfig{}, err
	}

	return t, nil
}

// ConnectionString returns connection string generated from configuration.
func (c DBConfig) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Db, c.SSLMode,
	)
}
