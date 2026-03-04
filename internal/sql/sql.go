package sql

import (
	"context"
	"net/url"
	"artifactor/pkg/config"

	"github.com/jackc/pgx/v5"
)


var Conn *pgx.Conn

func OpenConnection(cfg *config.PgsqlConfig) error {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.Username, cfg.Password),
		Host:   cfg.Addr,
		Path:   cfg.Database,
	}

	conn, err := pgx.Connect(context.Background(), u.String())
	if err != nil {
		return err
	}

	Conn = conn
	return nil
}
