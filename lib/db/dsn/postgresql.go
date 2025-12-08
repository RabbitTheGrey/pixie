package dsn

import (
	"strings"
)

type PostgresqlDsn struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
	SslMode  string
}

func (dsn *PostgresqlDsn) GetConnectionString() string {
	return strings.Join([]string{
		"user=" + dsn.User,
		"password=" + dsn.Password,
		"host=" + dsn.Host,
		"port=" + dsn.Port,
		"dbname=" + dsn.Db,
		"sslmode=" + dsn.SslMode,
	}, " ")
}
