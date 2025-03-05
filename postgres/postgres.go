package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Configuration struct {
	Host     string
	Port     string
	DBName   string
	User     string
	Password string
	SSLMode  string
}

func (c *Configuration) DSN() string {
	return "host=" + c.Host +
		" port=" + c.Port +
		" dbname=" + c.DBName +
		" user=" + c.User +
		" password=" + c.Password +
		" sslmode=" + c.SSLMode
}

func Open(c Configuration) (*sql.DB, error) {
	result, err := sql.Open("postgres", c.DSN())
	if err != nil {
		return nil, err
	}

	if err := result.Ping(); err != nil {
		return nil, err
	}

	return result, nil
}
