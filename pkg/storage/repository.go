package storage

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Db *sql.DB
}

func NewPostgres() (Repository, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't open postgres")
		return Repository{}, fmt.Errorf("Open postgres connection: %w", err)
	}
	err = db.Ping()
	if err != nil {
		log.Error().Err(err).Msgf("Couldn't ping postgres")
		return Repository{}, fmt.Errorf("Ping postgres connection: %w", err)
	}
	log.Info().Msgf("Connected to postgres")
	return Repository{
		Db: db,
	}, nil
}
