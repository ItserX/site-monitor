package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sites (
			id UUID PRIMARY KEY,
			url TEXT NOT NULL,
			active BOOLEAN NOT NULL DEFAULT true
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) AddSite(ctx context.Context, site Site) (string, error) {
	if site.ID == "" {
		site.ID = uuid.New().String()
	}
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO sites (id, url, active) VALUES ($1, $2, $3)`,
		site.ID, site.URL, site.Active,
	)
	if err != nil {
		return "", err
	}
	return site.ID, nil
}

func (p *PostgresStorage) GetSites(ctx context.Context) ([]Site, error) {
	rows, err := p.db.QueryContext(ctx, `SELECT id, url, active FROM sites`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sites []Site
	for rows.Next() {
		var s Site
		if err := rows.Scan(&s.ID, &s.URL, &s.Active); err != nil {
			return nil, err
		}
		sites = append(sites, s)
	}
	return sites, nil
}

func (p *PostgresStorage) GetSiteByID(ctx context.Context, id string) (*Site, error) {
	var s Site
	err := p.db.QueryRowContext(ctx,
		`SELECT id, url, active FROM sites WHERE id=$1`, id,
	).Scan(&s.ID, &s.URL, &s.Active)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (p *PostgresStorage) UpdateSite(ctx context.Context, site Site) error {
	_, err := p.db.ExecContext(ctx,
		`UPDATE sites SET url=$1, active=$2 WHERE id=$3`,
		site.URL, site.Active, site.ID,
	)
	return err
}

func (p *PostgresStorage) DeleteSite(ctx context.Context, id string) error {
	_, err := p.db.ExecContext(ctx,
		`DELETE FROM sites WHERE id=$1`, id,
	)
	return err
}
