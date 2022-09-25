package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type snake interface {
	info() BattlesnakeInfoResponse
	start(state GameState)
	end(state GameState)
	move(state GameState) BattlesnakeMoveResponse
}

// SnakeServer Serves a snake for battle.
type SnakeServer struct {
	snake snake
}

func NewSnakeServer(snake snake) *SnakeServer {
	return &SnakeServer{
		snake: snake,
	}
}

func (s *SnakeServer) HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := s.snake.info()
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func (s *SnakeServer) HandleStart(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}
	s.snake.start(state)
}

func (s *SnakeServer) HandleMove(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}
	response := s.snake.move(state)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

func (s *SnakeServer) HandleEnd(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}
	s.snake.end(state)
}

// Middleware

const ServerID = "battlesnake/github/starter-snake-go"

func RunServer(snake snake, port string) {
	server := NewSnakeServer(snake)
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Server", ServerID)
		server.HandleIndex(writer, request)
	})
	http.HandleFunc("/start", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Server", ServerID)
		server.HandleStart(writer, request)
	})
	http.HandleFunc("/move", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Server", ServerID)
		server.HandleMove(writer, request)
	})
	http.HandleFunc("/end", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Server", ServerID)
		server.HandleEnd(writer, request)
	})
	log.Printf("Running Battlesnake at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
