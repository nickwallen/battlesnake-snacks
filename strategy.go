package main

import (
	"github.com/rs/zerolog/log"
	"math"
)

// StayInBounds A strategy that kee[s the snake within the game boundaries.
type StayInBounds struct {
}

func (s *StayInBounds) move(state GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	if head.Right().X >= state.Board.Width || head.Right().Y < 0 {
		scorecard.Unsafe(RIGHT)
	}
	if head.Left().X < 0 || head.Left().X >= state.Board.Width {
		scorecard.Unsafe(LEFT)
	}
	if head.Up().Y >= state.Board.Height || head.Up().Y < 0 {
		scorecard.Unsafe(UP)
	}
	if head.Down().Y < 0 || head.Down().Y >= state.Board.Height {
		scorecard.Unsafe(DOWN)
	}
}

// NoCollisions allows a snake to avoid collisions with other snakes and itself.
type NoCollisions struct {
}

func (a *NoCollisions) move(state GameState, scorecard *Scorecard) {
	// Avoid other snakes and the snake itself
	var avoid []Body
	for _, opponent := range state.Board.Snakes {
		avoid = append(avoid, opponent.Body)
	}
	avoid = append(avoid, state.You.Body)

	head := headOfSnake(state)
	for _, body := range avoid {
		for _, coord := range body {
			if head.Right() == coord {
				scorecard.Unsafe(RIGHT)
			}
			if head.Left() == coord {
				scorecard.Unsafe(LEFT)
			}
			if head.Up() == coord {
				scorecard.Unsafe(UP)
			}
			if head.Down() == coord {
				scorecard.Unsafe(DOWN)
			}
		}
	}
}

// MoveToFood allows a snake to move toward the nearest food source.
type MoveToFood struct {
	weight Score
}

func (m MoveToFood) move(state GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	closestFood, err := findNearbyFood(state, head)
	if err != nil {
		return
	}

	// Incentivize moves that take us closer to the food
	if head.X > closestFood.X {
		scorecard.Add(LEFT, m.weight)
	}
	if head.X < closestFood.X {
		scorecard.Add(RIGHT, m.weight)
	}
	if head.Y > closestFood.Y {
		scorecard.Add(DOWN, m.weight)
	}
	if head.Y < closestFood.Y {
		scorecard.Add(UP, m.weight)
	}
}

type NoFoodFound struct {
}

func (n NoFoodFound) Error() string {
	return "No food found"
}

func findNearbyFood(state GameState, head Coord) (Coord, error) {
	var closestFood Coord
	if len(state.Board.Food) == 0 {
		return closestFood, NoFoodFound{}
	}
	minDist := math.MaxInt
	for _, food := range state.Board.Food {
		dist := head.DistanceTo(food)
		if dist < minDist {
			minDist = dist
			closestFood = food
		}
	}
	log.Debug().Msgf("Found closest food: %s", closestFood)
	return closestFood, nil
}

// MoveToCenter snakes should prefer moving toward the center.
type MoveToCenter struct {
	weight Score
}

func (m *MoveToCenter) move(state GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	centerX := float32(state.Board.Width) / float32(2)
	offsetX := float32(head.X) - centerX
	if offsetX < 0 {
		scorecard.Add(RIGHT, m.weight)
	}
	if offsetX > 0 {
		scorecard.Add(LEFT, m.weight)
	}
	centerY := float32(state.Board.Height) / float32(2)
	offsetY := float32(head.Y) - centerY
	if offsetY < 0 {
		scorecard.Add(UP, m.weight)
	}
	if offsetY > 0 {
		scorecard.Add(DOWN, m.weight)
	}
}

// MoveFromBiggerSnakes allows a snake to move away from larger snakes
type MoveFromBiggerSnakes struct {
	weight Score
}

func (m MoveFromBiggerSnakes) move(state GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	biggerSnake := findNearbySnake(state, head)

	// Penalize moves that take us closer to the bigger snake
	if head.X > biggerSnake.X {
		scorecard.Add(LEFT, -m.weight)
	} else {
		scorecard.Add(RIGHT, -m.weight)
	}
	if head.Y > biggerSnake.Y {
		scorecard.Add(DOWN, -m.weight)
	} else {
		scorecard.Add(UP, -m.weight)
	}
}

func findNearbySnake(state GameState, head Coord) Coord {
	var closestSnake Coord
	minDist := math.MaxInt
	for _, snake := range state.Board.Snakes {
		if snake.Length >= state.You.Length {
			dist := head.DistanceTo(snake.Head)
			if dist < minDist {
				minDist = dist
				closestSnake = snake.Head
			}
		}
	}
	return closestSnake
}
