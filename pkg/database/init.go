package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type StartUpOptions struct {
	DBHost     string
	DBPort     int
	DBName     string
	DBUsername string
	DBPassword string
}

func StartDbStore(opts StartUpOptions) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		opts.DBUsername, opts.DBPassword, opts.DBHost, opts.DBPort, opts.DBName))
	if err != nil {
		return nil, err
	}
	return db, nil
}
