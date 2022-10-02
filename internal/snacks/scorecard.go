package snacks

import (
	"fmt"
	b "github.com/nickwallen/battlesnake-snacks/internal/battlesnake"
	"math"
)

type Score int

func (s Score) String() string {
	return fmt.Sprintf("%d", s)
}

type Scorecard struct {
	state       b.GameState
	moves       map[b.Move]Score
	defaultMove b.Move // the default move is played if there are no safe moves
}

func NewScorecard(state b.GameState) *Scorecard {
	return &Scorecard{
		state: state,
		moves: map[b.Move]Score{
			b.LEFT:  0,
			b.RIGHT: 0,
			b.UP:    0,
			b.DOWN:  0,
		},
		defaultMove: b.DOWN,
	}
}

func (s *Scorecard) Add(move b.Move, toAdd Score) Score {
	if current, ok := s.moves[move]; ok {
		s.moves[move] = current + toAdd
		return current + toAdd
	}
	// Ignore an attempt to add to a move that was already marked unsafe
	return Score(0)
}

// Unsafe Marks a move as unsafe.
func (s *Scorecard) Unsafe(move b.Move) {
	delete(s.moves, move)
}

// Best Returns the move with the best score.
func (s *Scorecard) Best() b.Move {
	if len(s.moves) == 0 {
		logger(s.state).Msg("No safe moves!")
		return s.defaultMove
	}

	bestScore := Score(math.MinInt)
	var bestMove b.Move
	for move, score := range s.moves {
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}
	logger(s.state).Msgf("Chose %s as best from %v", bestMove, s.moves)
	return bestMove
}

func (s *Scorecard) SafeMoves() []b.Move {
	moves := make([]b.Move, 0)
	for k := range s.moves {
		moves = append(moves, k)
	}
	return moves
}

func (s *Scorecard) Scores() map[b.Move]Score {
	moves := make(map[b.Move]Score)
	for k, v := range s.moves {
		moves[k] = v
	}
	return moves
}

type LoggingScorecard struct {
	label     string      // the label prefixed to all logging
	state     b.GameState // the game state
	scorecard *Scorecard  // the actual scorecard
}

func NewLoggingScorecard(label string, state b.GameState, scorecard *Scorecard) *LoggingScorecard {
	return &LoggingScorecard{
		label:     label,
		state:     state,
		scorecard: scorecard,
	}
}

func (s *LoggingScorecard) Add(move b.Move, toAdd Score) Score {
	if toAdd != Score(0) {
		debug(s.state).Msgf("%s: %s +%d", s.label, move, toAdd)
	}
	return s.scorecard.Add(move, toAdd)
}

// Unsafe Marks a move as unsafe.
func (s *LoggingScorecard) Unsafe(move b.Move) {
	debug(s.state).Msgf("%s: %s 🚫", s.label, move)
	s.scorecard.Unsafe(move)
}

// Best Returns the move with the best score.
func (s *LoggingScorecard) Best() b.Move {
	return s.scorecard.Best()
}

func (s *LoggingScorecard) SafeMoves() []b.Move {
	return s.scorecard.SafeMoves()
}

func (s *LoggingScorecard) Scores() map[b.Move]Score {
	return s.scorecard.Scores()
}
