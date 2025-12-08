package dsn

import (
	"strings"
)

type MysqlDsn struct {
	User     string
	Password string
	Host     string
	Port     string
	Db       string
}

func (dsn *MysqlDsn) GetConnectionString() string {
	slice := []string{
		dsn.User, ":", dsn.Password, "@tcp(",
		dsn.Host, ":", dsn.Port, ")/", dsn.Db,
	}
	return strings.Join(slice, "")
}
