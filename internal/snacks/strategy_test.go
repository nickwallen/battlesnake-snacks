package snacks

import (
	"github.com/nickwallen/battlesnake-snacks/internal/battlesnake"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_StayInBounds_BottomLeft(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{0, 0},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []battlesnake.Move{battlesnake.RIGHT, battlesnake.UP}, scorecard.SafeMoves())
}

func Test_StayInBounds_BottomRight(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{4, 0},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []battlesnake.Move{battlesnake.LEFT, battlesnake.UP}, scorecard.SafeMoves())
}

func Test_StayInBounds_TopRight(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []battlesnake.Move{battlesnake.DOWN, battlesnake.LEFT}, scorecard.SafeMoves())
}

func Test_StayInBounds_TopLeft(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{0, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []battlesnake.Move{battlesnake.RIGHT, battlesnake.DOWN}, scorecard.SafeMoves())
}

func Test_StayInBounds_Middle(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.ElementsMatch(t, []battlesnake.Move{battlesnake.RIGHT, battlesnake.LEFT, battlesnake.UP, battlesnake.DOWN}, scorecard.SafeMoves())
}

func Test_StayInBounds_OutOfBounds(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{-5, -5},
		},
	}
	scorecard := NewScorecard(state)
	strategy := StayInBounds{}
	strategy.move(state, scorecard)
	require.Equal(t, []battlesnake.Move{}, scorecard.SafeMoves())
}

func Test_NoCollision_AvoidSelf(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{1, 1},
			Body: []battlesnake.Coord{
				{1, 1},
				{2, 1},
			},
		},
	}
	scorecard := NewScorecard(state)
	strategy := NoCollisions{}
	strategy.move(state, scorecard)
	require.NotContains(t, scorecard.SafeMoves(), battlesnake.RIGHT)
}

func Test_NoCollision_AvoidOpponent(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					Head: battlesnake.Coord{0, 1},
					Body: []battlesnake.Coord{
						{0, 1},
					},
					Length: 1,
				},
			},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{1, 1},
			Body: []battlesnake.Coord{
				{1, 1},
			},
			Length: 1,
		},
	}
	scorecard := NewScorecard(state)
	strategy := NoCollisions{}
	strategy.move(state, scorecard)
	require.NotContains(t, scorecard.SafeMoves(), battlesnake.LEFT)
}

func Test_MoveToCenter_BottomLeft(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{1, 1},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.UP])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
}

func Test_MoveToCenter_BottomRight(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{4, 1},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.UP])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
}

func Test_MoveToCenter_TopRight(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(2), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToCenter_Center(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 6,
			Width:  6,
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{3, 3},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToCenter{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToClosestFood_NoFood(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food:   []battlesnake.Coord{},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{4, 4},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToClosestFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToClosestFood_Nearest(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food: []battlesnake.Coord{
				{1, 1},
				{4, 4},
			},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToClosestFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(10), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(10), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToClosestFood_MoveUpOnly(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food: []battlesnake.Coord{
				{2, 3},
			},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToClosestFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(10), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToClosestFood_MoveRightOnly(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food: []battlesnake.Coord{
				{3, 2},
			},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToClosestFood{weight: Score(10)}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(10), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_AvoidBiggerSnakes_BiggerSnake(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "bigger-snake",
					Head:   battlesnake.Coord{3, 3},
					Length: 5,
				},
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AvoidBiggerSnakes{weight: 1.0}
	strategy.move(state, scorecard)
	require.Equal(t, Score(6), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(6), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_AvoidBiggerSnakes_SmallerSnake(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "smaller-snake",
					Head:   battlesnake.Coord{3, 3},
					Length: 1,
				},
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AvoidBiggerSnakes{weight: 1.0}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_AvoidBiggerSnakes_NoSnakes(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AvoidBiggerSnakes{weight: 1.0}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_Board_IsEmpty(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 2,
			Width:  2,
			Snakes: []battlesnake.Snake{
				{
					ID:   "you",
					Head: battlesnake.Coord{0, 0},
					Body: []battlesnake.Coord{
						{0, 0},
						{1, 1},
					},
					Length: 2,
				},
			},
			Hazards: []battlesnake.Coord{
				{0, 1},
			},
		},
	}
	board := NewBoard(state)
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{0, 0}))
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{0, 1}))
	require.Equal(t, true, board.isEmpty(battlesnake.Coord{1, 0}))
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{1, 1}))
}

func Test_Board_OutOfBounds(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 2,
			Width:  2,
		},
	}
	board := NewBoard(state)
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{-1, 0}))
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{2, 2}))
	require.Equal(t, false, board.isEmpty(battlesnake.Coord{-1, -2}))
}

func Test_MoveToSpace(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 3,
			Width:  3,
			Snakes: []battlesnake.Snake{
				{
					ID:   "you",
					Head: battlesnake.Coord{2, 1},
					Body: []battlesnake.Coord{
						{2, 1},
						{1, 1},
						{1, 0},
						{0, 0},
					},
					Length: 4,
				},
			},
		},
		You: battlesnake.Snake{
			ID:   "you",
			Head: battlesnake.Coord{2, 1},
			Body: []battlesnake.Coord{
				{2, 1},
				{1, 1},
				{1, 0},
				{0, 0},
			},
			Length: 4,
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToSpace{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(6), scorecard.Scores()[battlesnake.UP])
	require.Equal(t, Score(1), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
}

func Test_MoveToFood_NoFood(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food:   []battlesnake.Coord{},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_MoveToFood_WithFood(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Food: []battlesnake.Coord{
				{1, 1},
				{4, 4},
			},
		},
		You: battlesnake.Snake{
			Head: battlesnake.Coord{2, 2},
		},
	}
	scorecard := NewScorecard(state)
	strategy := MoveToFood{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(9), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(9), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(6), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(6), scorecard.Scores()[battlesnake.UP])
}

func Test_AvoidDeadEnds_DeadEnd(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:   "you",
					Head: battlesnake.Coord{2, 1},
					Body: []battlesnake.Coord{
						{2, 1},
						{1, 1},
						{0, 1},
					},
					Length: 3,
				},
			},
			Hazards: []battlesnake.Coord{
				{3, 0},
				{3, 2},
				{4, 0},
				{4, 2},
			},
		},
		You: battlesnake.Snake{
			ID:   "you",
			Head: battlesnake.Coord{2, 1},
			Body: []battlesnake.Coord{
				{2, 1},
				{1, 1},
				{0, 1},
			},
			Length: 3,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AvoidDeadEnds{}
	strategy.move(state, scorecard)

	// Right is a dead-end!
	require.NotContains(t, scorecard.SafeMoves(), battlesnake.RIGHT)
}

func Test_AvoidDeadEnds_NotADeadEnd(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:   "you",
					Head: battlesnake.Coord{2, 1},
					Body: []battlesnake.Coord{
						{1, 1},
						{0, 1},
					},
					Length: 2,
				},
			},
			Hazards: []battlesnake.Coord{
				{3, 0},
				{3, 2},
				{4, 0},
				{4, 2},
			},
		},
		You: battlesnake.Snake{
			ID:   "you",
			Head: battlesnake.Coord{2, 1},
			Body: []battlesnake.Coord{
				{2, 1},
				{1, 1},
			},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AvoidDeadEnds{}
	strategy.move(state, scorecard)
	require.Contains(t, scorecard.SafeMoves(), battlesnake.RIGHT)
}

func Test_AttackSmallerSnakes_SmallerSnake(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "smaller-snake",
					Head:   battlesnake.Coord{3, 3},
					Length: 2,
				},
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 3,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 3,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AttackSmallerSnakes{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(9), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(9), scorecard.Scores()[battlesnake.UP])
}

func Test_AttackSmallerSnakes_BiggerSnake(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "bigger-snake",
					Head:   battlesnake.Coord{3, 3},
					Length: 5,
				},
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AttackSmallerSnakes{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_AttackSmallerSnakes_EqualSnake(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "equal-size-snake",
					Head:   battlesnake.Coord{3, 3},
					Length: 2,
				},
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AttackSmallerSnakes{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}

func Test_AttackSmallerSnakes_NoSnakes(t *testing.T) {
	state := battlesnake.GameState{
		Board: battlesnake.Board{
			Height: 5,
			Width:  5,
			Snakes: []battlesnake.Snake{
				{
					ID:     "you",
					Head:   battlesnake.Coord{2, 2},
					Length: 2,
				},
			},
		},
		You: battlesnake.Snake{
			ID:     "you",
			Head:   battlesnake.Coord{2, 2},
			Length: 2,
		},
	}
	scorecard := NewScorecard(state)
	strategy := AttackSmallerSnakes{weight: 1.5}
	strategy.move(state, scorecard)
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.LEFT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.DOWN])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.RIGHT])
	require.Equal(t, Score(0), scorecard.Scores()[battlesnake.UP])
}
