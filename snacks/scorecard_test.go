package snacks

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func state() GameState {
	return GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{0, 0},
		},
	}
}

func Test_Scorecard_Add(t *testing.T) {
	s := NewScorecard(state())
	require.Equal(t, Score(10), s.Add(LEFT, 10))
	require.Equal(t, Score(20), s.Add(LEFT, 10))
	require.Equal(t, Score(30), s.Add(LEFT, 10))
}

func Test_Scorecard_Best(t *testing.T) {
	s := NewScorecard(state())
	s.Add(LEFT, 10)
	s.Add(RIGHT, 5)
	s.Add(UP, 5)
	s.Add(DOWN, 5)
	require.Equal(t, LEFT, s.Best())
}

func Test_Scorecard_Unsafe(t *testing.T) {
	s := NewScorecard(state())
	s.Unsafe(LEFT)
	s.Add(LEFT, 10)
	s.Add(RIGHT, 5)
	require.Equal(t, RIGHT, s.Best())
}

func Test_Scorecard_Scores(t *testing.T) {
	s := NewScorecard(state())
	safeMoves := s.Scores()
	require.Equal(t, Score(0), safeMoves[LEFT])
	require.Equal(t, Score(0), safeMoves[RIGHT])
	require.Equal(t, Score(0), safeMoves[UP])
	require.Equal(t, Score(0), safeMoves[DOWN])

	s.Add(LEFT, 1)
	s.Add(RIGHT, 2)
	s.Add(UP, 3)
	s.Add(DOWN, 4)
	safeMoves = s.Scores()
	require.Equal(t, Score(1), safeMoves[LEFT])
	require.Equal(t, Score(2), safeMoves[RIGHT])
	require.Equal(t, Score(3), safeMoves[UP])
	require.Equal(t, Score(4), safeMoves[DOWN])
}
