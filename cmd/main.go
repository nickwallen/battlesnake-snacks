package main

import (
	"github.com/nickwallen/battlesnake-snacks/snacks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	EnvPort  = "PORT"
	EnvSnake = "SNAKE"

	PortDefault = "8000"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Which port to use?
	port := os.Getenv(EnvPort)
	if len(port) == 0 {
		port = PortDefault
	}

	// Which snake will battle?
	var snake snake
	switch os.Getenv(EnvSnake) {
	case "DUMB":
		snake = snacks.NewDumbSnake()
	case "HUNGRY":
		snake = snacks.NewHungrySnake()
	case "LATEST":
		snake = snacks.NewNextGenSnake()
	default:
		log.Fatal().Msgf("Env var '%s' is missing.", EnvSnake)
	}

	RunServer(snake, port)
}
