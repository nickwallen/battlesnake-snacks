package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Which port to use?
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	// Which snake will battle?
	var snake snake
	switch os.Getenv("SNAKE") {
	case "DUMB":
		snake = NewDumbSnake()
	case "HUNGRY":
		snake = NewHungrySnake()
	default:
		snake = NewNextGenSnake()
	}

	RunServer(snake, port)
}
