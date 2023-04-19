package postgres

type Config struct {
	DSN            string `mapstructure:"dsn"`
	MaxConnections int32  `mapstructure:"max_connections"`
}
