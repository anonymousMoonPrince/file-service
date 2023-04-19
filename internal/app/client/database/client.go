package database

import "github.com/jackc/pgtype/pgxtype"

type Client interface {
	pgxtype.Querier
}
