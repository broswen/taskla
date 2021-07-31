package main

import (
	"github.com/broswen/taskla/cmd/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("Starting Taskla")
	server, err := server.New()
	if err != nil {
		log.Fatal().Err(err).Msg("Couldn't initialize server")
		return
	}

	log.Info().Msg("Setting routes")
	server.Routes()

	if err := server.Start(); err != nil {
		log.Fatal().Err(err).Msg("Couldn't start server")
	}
}
