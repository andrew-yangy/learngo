package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ddvkid/learngo/internal/config"
	"github.com/ddvkid/learngo/internal/storage"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
	"net/url"
	"time"
)

type Tx struct {
	*sqlx.Tx
	Source
}

type PgStore struct {
	*sqlx.DB
	Source
}

type Source struct {
	Queryable
}

type Queryable interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryxContext(context.Context, string, ...interface{}) (*sqlx.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	SelectContext(context.Context, interface{}, string, ...interface{}) error
	GetContext(context.Context, interface{}, string, ...interface{}) error
	NamedExecContext(context.Context, string, interface{}) (sql.Result, error)
}

func pgMaxConnectionLifetime() time.Duration {
	return time.Minute * 10
}

func NewStore(connString string) (*PgStore, error) {
	sqldb, err := sql.Open("nrpostgres", connString)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(sqldb, "postgres")
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(pgMaxConnectionLifetime())
	return &PgStore{DB: db, Source: Source{Queryable: db}}, nil
}

func GetUrl(host, port, db, username, password string) *url.URL {
	pgURL := &url.URL{
		Host:   fmt.Sprintf("%s:%s", host, port),
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Path:   db,
	}
	q := pgURL.Query()
	pgURL.RawQuery = q.Encode()
	return pgURL
}

func GetUrlByStage() *url.URL {
	return GetUrl(config.PGEndpoint, config.PGPort, config.PGDbName, config.PGUserName, config.PGPassword)
}

func (s PgStore) BeginTx(ctx context.Context) (storage.Tx, error) {
	tx, err := s.DB.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	return Tx{Tx: tx, Source: Source{Queryable: tx}}, err
}

func (s PgStore) Ping() error {
	return s.DB.Ping()
}
