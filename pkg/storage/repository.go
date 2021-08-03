package storage

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	db *sql.DB
}

func New() Repository {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msgf("connect: Couldn't open postgres")
	}
	log.Info().Msgf("Connected to postgres")
	return Repository{db}
}

func (r *Repository) DB() *sql.DB {
	err := r.db.Ping()
	for err != nil {
		log.Error().Err(err).Msg("DB: Couldn't ping postgres")
		time.Sleep(5 * time.Second)
		err = r.db.Ping()
	}

	return r.db
}
