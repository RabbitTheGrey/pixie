package dsn

type SqliteDsn struct {
	Path string
}

func (dsn *SqliteDsn) GetConnectionString() string {
	return dsn.Path
}
