package config

type dbType string

const (
	Sqlite dbType = "sqlite3"
)

type DB struct {
	Type dbType
	Dsn  string
}
