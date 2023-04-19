package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

const connectTimeout = time.Minute * 5

type Client struct {
	*pgxpool.Pool
}

func NewClient(config Config) *Client {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		logrus.WithError(err).Fatal("parse postgres dsn failed")
	}
	poolConfig.MaxConns = config.MaxConnections

	pgxPoolConn, err := pgxpool.ConnectConfig(ctx, poolConfig)
	if err != nil {
		logrus.WithError(err).Fatal("connect to postgres failed")
	}

	return &Client{
		pgxPoolConn,
	}
}
