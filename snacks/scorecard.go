package snacks

import (
	"fmt"
	"math"
)

type Score int

func (s Score) String() string {
	return fmt.Sprintf("%d", s)
}

type Scorecard struct {
	state       GameState
	moves       map[Move]Score
	defaultMove Move // the default move is played if there are no safe moves
}

func NewScorecard(state GameState) *Scorecard {
	return &Scorecard{
		state: state,
		moves: map[Move]Score{
			LEFT:  0,
			RIGHT: 0,
			UP:    0,
			DOWN:  0,
		},
		defaultMove: DOWN,
	}
}

func (s *Scorecard) Add(move Move, toAdd Score) Score {
	if current, ok := s.moves[move]; ok {
		s.moves[move] = current + toAdd
		return current + toAdd
	}
	// Ignore an attempt to add to a move that was already marked unsafe
	return Score(0)
}

// Unsafe Marks a move as unsafe.
func (s *Scorecard) Unsafe(move Move) {
	delete(s.moves, move)
}

// Best Returns the move with the best score.
func (s *Scorecard) Best() Move {
	if len(s.moves) == 0 {
		logger(s.state).Msg("No safe moves!")
		return s.defaultMove
	}

	bestScore := Score(math.MinInt)
	var bestMove Move
	for move, score := range s.moves {
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	logger(s.state).Msgf("Chose best move from %v", s.moves)
	return bestMove
}

func (s *Scorecard) SafeMoves() []Move {
	moves := []Move{}
	for k := range s.moves {
		moves = append(moves, k)
	}
	return moves
}

func (s *Scorecard) Scores() map[Move]Score {
	moves := make(map[Move]Score)
	for k, v := range s.moves {
		moves[k] = v
	}
	return moves
}
