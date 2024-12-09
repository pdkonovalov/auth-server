package postgres

import (
	"context"
	"errors"

	"github.com/pdkonovalov/auth-server/pkg/config"
	"github.com/pdkonovalov/auth-server/pkg/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	pool *pgxpool.Pool
}

func Init(config *config.Config) (storage.Storage, error) {
	pool, err := pgxpool.New(context.Background(), config.DatabaseUrl)
	if err != nil {
		return nil, err
	}
	_, err = pool.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS valid_jwt (
		guid UUID NOT NULL,
		jti UUID NOT NULL);
		CREATE INDEX IF NOT EXISTS jti_alias ON valid_jwt(jti);
	`)
	if err != nil {
		return nil, err
	}
	return &postgres{pool}, nil
}

func (db *postgres) Shutdown() error {
	db.pool.Close()
	return nil
}

func (db *postgres) WriteNewJti(guid string) (string, error) {
	jti := uuid.New().String()
	_, err := db.pool.Exec(context.Background(),
		`INSERT INTO valid_jwt(guid, jti) VALUES($1, $2);`, guid, jti)
	if err != nil {
		return "", err
	}
	return jti, nil
}

func (db *postgres) FindJti(jti string) (string, bool, error) {
	var guid string
	err := db.pool.QueryRow(context.Background(),
		`SELECT guid FROM valid_jwt WHERE jti = $1`, jti).Scan(&guid)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return guid, true, nil
}

func (db *postgres) DeleteJti(jti string) error {
	_, err := db.pool.Exec(context.Background(),
		`DELETE FROM valid_jwt WHERE jti = $1`, jti)
	return err
}
