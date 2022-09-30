package snacks

import (
	b "github.com/nickwallen/battlesnake-snacks/battlesnake"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type strategy interface {
	move(state b.GameState, scorecard *Scorecard)
}

type StrategyDrivenSnake struct {
	name       string
	author     string
	color      string
	head       string
	tail       string
	strategies []strategy
}

func NewDumbSnake() *StrategyDrivenSnake {
	return &StrategyDrivenSnake{
		name:   "dumb",
		author: "nickwallen",
		color:  "#b5ca60",
		head:   "dead",
		tail:   "do-sammy",
		strategies: []strategy{
			&StayInBounds{},
			&NoCollisions{},
		},
	}
}

func NewHungrySnake() *StrategyDrivenSnake {
	return &StrategyDrivenSnake{
		name:   "hungry",
		author: "nickwallen",
		color:  "#2F4538",
		head:   "ski",
		tail:   "coffee",
		strategies: []strategy{
			&StayInBounds{},
			&NoCollisions{},
			&MoveToClosestFood{weight: 20},
			&MoveToCenter{weight: 10},
		},
	}
}

func NewLatestSnake() *StrategyDrivenSnake {
	return &StrategyDrivenSnake{
		name:   "next-gen",
		author: "nickwallen",
		color:  "#00BB2D",
		head:   "regular",
		tail:   "regular",
		strategies: []strategy{
			&StayInBounds{},
			&NoCollisions{},
			&MoveToClosestFood{weight: 5},
			//&MoveToFood{weight: 1},
			&MoveToCenter{weight: 4},
			&AvoidBiggerSnakes{weight: 1.0},
			&AvoidDeadEnds{weight: 1},
		},
	}
}

// Info is called when you create your Battlesnake on play.b.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func (s *StrategyDrivenSnake) Info() b.InfoResponse {
	return b.InfoResponse{
		APIVersion: "1",
		Author:     s.author,
		Color:      s.color,
		Head:       s.head,
		Tail:       s.tail,
	}
}

// Start is called when your Battlesnake begins a game
func (s *StrategyDrivenSnake) Start(state b.GameState) {
	logger(state).
		Str("snake", s.name).
		Msg("start")
}

// End is called when your Battlesnake finishes a game
func (s *StrategyDrivenSnake) End(state b.GameState) {
	var gameResult string
	isDraw := len(state.Board.Snakes) == 0
	if isDraw {
		gameResult = "Draw"
	} else {
		isWinner := state.You.ID == state.Board.Snakes[0].ID
		if isWinner {
			gameResult = "Won"
		} else {
			gameResult = "Lost"
		}
	}
	logger(state).
		Str("snake", s.name).
		Msgf("%s in %d move(s)", gameResult, state.Turn+1)
}

// Move is called on every turn and returns your next move
// Valid moves are UP, DOWN, LEFT, or RIGHT
// See https://docs.b.com/api/example-move for available data
func (s *StrategyDrivenSnake) Move(state b.GameState) b.MoveResponse {
	scorecard := NewScorecard(state)
	for _, strategy := range s.strategies {
		strategy.move(state, scorecard)
	}
	move := scorecard.Best()
	logger(state).Stringer("move", move).Msg("moved")
	return b.MoveResponse{Move: move}
}

func logger(state b.GameState) *zerolog.Event {
	event := log.Info().
		Str("game-id", state.Game.ID).
		Int("turn", state.Turn).
		Stringer("head", headOfSnake(state)).
		Int("health", state.You.Health).
		Int("length", state.You.Length)
	return event
}

func debug(state b.GameState) *zerolog.Event {
	event := log.Debug().
		Str("game-id", state.Game.ID).
		Int("turn", state.Turn).
		Stringer("head", headOfSnake(state))
	return event
}

// head Returns the coordinates of the snake's head.
func headOfSnake(state b.GameState) b.Coord {
	return state.You.Head
}