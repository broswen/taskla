package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Db *sql.DB
}

func NewPostgres() (Repository, error) {
	connStr := "user=taskla password=password host=postgres dbname=taskla port=5432 sslmode=disable"
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
