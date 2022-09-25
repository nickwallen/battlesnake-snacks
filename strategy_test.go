package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_StayInBounds_BottomLeft(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{0, 0},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []Move{RIGHT, UP}, scorecard.SafeMoves())
}

func Test_StayInBounds_BottomRight(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{4, 0},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []Move{LEFT, UP}, scorecard.SafeMoves())
}

func Test_StayInBounds_TopRight(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []Move{DOWN, LEFT}, scorecard.SafeMoves())
}

func Test_StayInBounds_TopLeft(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{0, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []Move{RIGHT, DOWN}, scorecard.SafeMoves())
}

func Test_StayInBounds_Middle(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []Move{RIGHT, LEFT, UP, DOWN}, scorecard.SafeMoves())
}

func Test_StayInBounds_OutOfBounds(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{-5, -5},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.Equal(t, []Move{}, scorecard.SafeMoves())
}

func Test_NoCollision_AvoidSelf(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{1, 1},
			Body: []Coord{
				{1, 1},
				{2, 1},
			},
		},
	}
	scorecard := NewScorecard(state)
	strategy := NoCollisions{}
	strategy.move(state, scorecard)
	require.NotContains(t, scorecard.SafeMoves(), RIGHT)
}

func Test_NoCollision_AvoidOpponent(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
			Snakes: []Battlesnake{
				{
					Head: Coord{0, 1},
					Body: []Coord{
						{0, 1},
					},
					Length: 1,
				},
			},
		},
		You: Battlesnake{
			Head: Coord{1, 1},
			Body: []Coord{
				{1, 1},
			},
			Length: 1,
		},
	}
	scorecard := NewScorecard(state)
	strategy := NoCollisions{}
	strategy.move(state, scorecard)
	require.NotContains(t, scorecard.SafeMoves(), LEFT)
}

func Test_MoveToCenter_BottomLeft(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{1, 1},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(10), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(10), scorecard.Scores()[UP])
	require.Equal(t, Score(0), scorecard.Scores()[LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
}

func Test_MoveToCenter_BottomRight(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{4, 1},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(10), scorecard.Scores()[LEFT])
	require.Equal(t, Score(10), scorecard.Scores()[UP])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
}

func Test_MoveToCenter_TopRight(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
		},
		You: Battlesnake{
			Head: Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(10), scorecard.Scores()[LEFT])
	require.Equal(t, Score(10), scorecard.Scores()[DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[UP])
}

func Test_MoveToCenter_Center(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 6,
			Width:  6,
		},
		You: Battlesnake{
			Head: Coord{3, 3},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[UP])
}

func Test_MoveToFood_NoFood(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
			Food:   []Coord{},
		},
		You: Battlesnake{
			Head: Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[UP])
}

func Test_MoveToFood_Nearest(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
			Food: []Coord{
				{1, 1},
				{4, 4},
			},
		},
		You: Battlesnake{
			Head: Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(10), scorecard.Scores()[LEFT])
	require.Equal(t, Score(10), scorecard.Scores()[DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[UP])
}

func Test_MoveToFood_MoveUpOnly(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
			Food: []Coord{
				{2, 3},
			},
		},
		You: Battlesnake{
			Head: Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(10), scorecard.Scores()[UP])
}

func Test_MoveToFood_MoveRightOnly(t *testing.T) {
	state := GameState{
		Board: Board{
			Height: 5,
			Width:  5,
			Food: []Coord{
				{3, 2},
			},
		},
		You: Battlesnake{
			Head: Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[DOWN])
	require.Equal(t, Score(10), scorecard.Scores()[RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[UP])
}
