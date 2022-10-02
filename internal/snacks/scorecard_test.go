package snacks

import (
	"github.com/nickwallen/battlesnake-snacks/internal/battlesnake"
	"github.com/stretchr/testify/require"
	"testing"
)

func state() battlesnake.GameState {
	return battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{0, 0},
		},
	}
}

func Test_Scorecard_Add(t *testing.T) {
	s := NewScorecard(state())
	require.Equal(t, Score(10), s.Add(battlesnake.LEFT, 10))
	require.Equal(t, Score(20), s.Add(battlesnake.LEFT, 10))
	require.Equal(t, Score(30), s.Add(battlesnake.LEFT, 10))
}

func Test_Scorecard_Best(t *testing.T) {
	s := NewScorecard(state())
	s.Add(battlesnake.LEFT, 10)
	s.Add(battlesnake.RIGHT, 5)
	s.Add(battlesnake.UP, 5)
	s.Add(battlesnake.DOWN, 5)
	require.Equal(t, battlesnake.LEFT, s.Best())
}

func Test_Scorecard_Unsafe(t *testing.T) {
	s := NewScorecard(state())
	s.Unsafe(battlesnake.LEFT)
	s.Add(battlesnake.LEFT, 10)
	s.Add(battlesnake.RIGHT, 5)
	require.Equal(t, battlesnake.RIGHT, s.Best())
}

func Test_Scorecard_Scores(t *testing.T) {
	s := NewScorecard(state())
	safeMoves := s.Scores()
	require.Equal(t, Score(0), safeMoves[battlesnake.LEFT])
	require.Equal(t, Score(0), safeMoves[battlesnake.RIGHT])
	require.Equal(t, Score(0), safeMoves[battlesnake.UP])
	require.Equal(t, Score(0), safeMoves[battlesnake.DOWN])

	s.Add(battlesnake.LEFT, 1)
	s.Add(battlesnake.RIGHT, 2)
	s.Add(battlesnake.UP, 3)
	s.Add(battlesnake.DOWN, 4)
	safeMoves = s.Scores()
	require.Equal(t, Score(1), safeMoves[battlesnake.LEFT])
	require.Equal(t, Score(2), safeMoves[battlesnake.RIGHT])
	require.Equal(t, Score(3), safeMoves[battlesnake.UP])
	require.Equal(t, Score(4), safeMoves[battlesnake.DOWN])
}
