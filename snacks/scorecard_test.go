package snacks

import (
	b "github.com/nickwallen/battlesnake-snacks/battlesnake"
	"github.com/stretchr/testify/require"
	"testing"
)

func state() b.GameState {
	return b.GameState{
		Board: b.Board{
			Height: 5,
			Width:  5,
		},
		You: b.Snake{
			Head: b.Coord{0, 0},
		},
	}
}

func Test_Scorecard_Add(t *testing.T) {
	s := NewScorecard(state())
	require.Equal(t, Score(10), s.Add(b.LEFT, 10))
	require.Equal(t, Score(20), s.Add(b.LEFT, 10))
	require.Equal(t, Score(30), s.Add(b.LEFT, 10))
}

func Test_Scorecard_Best(t *testing.T) {
	s := NewScorecard(state())
	s.Add(b.LEFT, 10)
	s.Add(b.RIGHT, 5)
	s.Add(b.UP, 5)
	s.Add(b.DOWN, 5)
	require.Equal(t, b.LEFT, s.Best())
}

func Test_Scorecard_Unsafe(t *testing.T) {
	s := NewScorecard(state())
	s.Unsafe(b.LEFT)
	s.Add(b.LEFT, 10)
	s.Add(b.RIGHT, 5)
	require.Equal(t, b.RIGHT, s.Best())
}

func Test_Scorecard_Scores(t *testing.T) {
	s := NewScorecard(state())
	safeMoves := s.Scores()
	require.Equal(t, Score(0), safeMoves[b.LEFT])
	require.Equal(t, Score(0), safeMoves[b.RIGHT])
	require.Equal(t, Score(0), safeMoves[b.UP])
	require.Equal(t, Score(0), safeMoves[b.DOWN])

	s.Add(b.LEFT, 1)
	s.Add(b.RIGHT, 2)
	s.Add(b.UP, 3)
	s.Add(b.DOWN, 4)
	safeMoves = s.Scores()
	require.Equal(t, Score(1), safeMoves[b.LEFT])
	require.Equal(t, Score(2), safeMoves[b.RIGHT])
	require.Equal(t, Score(3), safeMoves[b.UP])
	require.Equal(t, Score(4), safeMoves[b.DOWN])
}
