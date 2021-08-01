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
	return Repository{}
}

func (r *Repository) DB() *sql.DB {
	if r.db == nil {
		var err error
		r.db, err = connect()
		if err != nil {
			log.Error().Err(err).Msg("DB: Couldn't connect to postgres")
		}
	}

	err := r.db.Ping()
	for err != nil {
		log.Error().Err(err).Msg("DB: Couldn't ping postgres")
		if r.db != nil {
			r.db.Close()
		}
		time.Sleep(5 * time.Second)
		r.db, err = connect()
		if err == nil {
			err = r.db.Ping()
		}
	}

	log.Info().Msgf("Connected to postgres")
	return r.db
}

func connect() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s host=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error().Err(err).Msgf("connect: Couldn't open postgres")
		return db, err
	}
	return db, nil
}
