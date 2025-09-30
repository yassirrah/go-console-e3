package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	DB *sql.DB
}

func Open(dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn) // e.g. postgres://user:pass@db:5432/weather?sslmode=disable
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	s := &Store{DB: db}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error { return s.DB.Close() }

func (s *Store) migrate() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS weather_queries (
	id SERIAL PRIMARY KEY,
	city TEXT NOT NULL,
	units TEXT NOT NULL,
	temperature_c NUMERIC NOT NULL,
	humidity_pct INT NOT NULL,
	conditions TEXT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);`
	_, err := s.DB.Exec(ddl)
	return err
}

type WeatherRow struct {
	City        string
	Units       string
	Temperature float64
	Humidity    int
	Conditions  string
}

func (s *Store) InsertWeather(ctx context.Context, w WeatherRow) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO weather_queries (city, units, temperature_c, humidity_pct, conditions)
		 VALUES ($1, $2, $3, $4, $5)`,
		w.City, w.Units, w.Temperature, w.Humidity, w.Conditions)
	return err
}
