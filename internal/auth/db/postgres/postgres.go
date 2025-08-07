package postgres

import (
	"context"
	"database/sql"
	"time"

	dbContract "github.com/barantoraman/microgate/internal/auth/db/contract"
	"github.com/barantoraman/microgate/pkg/config"
	_ "github.com/lib/pq"
)

type pqConnection struct {
	db *sql.DB
}

func (p *pqConnection) Close() {
	p.db.Close()
}

func (p *pqConnection) DB() *sql.DB {
	return p.db
}

func NewPostgresConnection(cfg config.AuthServiceConfigurations) (dbContract.DBConnection, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return &pqConnection{
		db: db,
	}, nil
}
