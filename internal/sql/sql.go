package sql

import (
	"fmt"
	"context"
	"net/url"
	"artifactor/pkg/config"

	"github.com/jackc/pgx/v5"
)


var Conn *pgx.Conn = nil

func OpenConnection(cfg *config.PgsqlConfig) error {
	Conn = nil
	if cfg == nil {
		return fmt.Errorf("Missing pgsql config")
	}

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
