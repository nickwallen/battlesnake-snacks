package snacks

import (
	"errors"
	b "github.com/nickwallen/battlesnake-snacks/battlesnake"
	"github.com/rs/zerolog/log"
	"math"
)

var (
	ErrNoBiggerSnakes = errors.New("no bigger snakes found")
	ErrNoFood         = errors.New("no food found")
)

// StayInBounds A strategy that kee[s the snake within the game boundaries.
type StayInBounds struct {
}

func (s *StayInBounds) move(state b.GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	if head.Right().X >= state.Board.Width || head.Right().Y < 0 {
		scorecard.Unsafe(b.RIGHT)
	}
	if head.Left().X < 0 || head.Left().X >= state.Board.Width {
		scorecard.Unsafe(b.LEFT)
	}
	if head.Up().Y >= state.Board.Height || head.Up().Y < 0 {
		scorecard.Unsafe(b.UP)
	}
	if head.Down().Y < 0 || head.Down().Y >= state.Board.Height {
		scorecard.Unsafe(b.DOWN)
	}
}

// NoCollisions allows a snake to avoid collisions with other snakes and itself.
type NoCollisions struct {
}

func (a *NoCollisions) move(state b.GameState, scorecard *Scorecard) {
	// Avoid other snakes and the snake itself
	var avoid []b.Body
	for _, opponent := range state.Board.Snakes {
		avoid = append(avoid, opponent.Body)
	}
	avoid = append(avoid, state.You.Body)

	head := headOfSnake(state)
	for _, body := range avoid {
		for _, coord := range body {
			if head.Right() == coord {
				scorecard.Unsafe(b.RIGHT)
			}
			if head.Left() == coord {
				scorecard.Unsafe(b.LEFT)
			}
			if head.Up() == coord {
				scorecard.Unsafe(b.UP)
			}
			if head.Down() == coord {
				scorecard.Unsafe(b.DOWN)
			}
		}
	}
}

// MoveToClosestFood allows a snake to move toward the nearest food source.
type MoveToClosestFood struct {
	weight Score
}

func (m MoveToClosestFood) move(state b.GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	closestFood, err := findNearbyFood(state, head)
	if err == ErrNoFood {
		return // Nothing to do
	}

	// Incentivize moves that take us closer to the food
	if head.X > closestFood.X {
		scorecard.Add(b.LEFT, m.weight)
	}
	if head.X < closestFood.X {
		scorecard.Add(b.RIGHT, m.weight)
	}
	if head.Y > closestFood.Y {
		scorecard.Add(b.DOWN, m.weight)
	}
	if head.Y < closestFood.Y {
		scorecard.Add(b.UP, m.weight)
	}
}

func findNearbyFood(state b.GameState, head b.Coord) (b.Coord, error) {
	var closestFood b.Coord
	if len(state.Board.Food) == 0 {
		return closestFood, ErrNoFood
	}
	minDist := math.MaxInt
	for _, food := range state.Board.Food {
		dist := head.DistanceTo(food)
		if dist < minDist {
			minDist = dist
			closestFood = food
		}
	}
	debug(state).Msgf("Found closest food: %s", closestFood)
	return closestFood, nil
}

// MoveToCenter snakes should prefer moving toward the center.
type MoveToCenter struct {
	weight Score
}

func (m *MoveToCenter) move(state b.GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	centerX := float32(state.Board.Width) / float32(2)
	offsetX := float32(head.X) - centerX
	if offsetX < 0 {
		scorecard.Add(b.RIGHT, m.weight)
	}
	if offsetX > 0 {
		scorecard.Add(b.LEFT, m.weight)
	}
	centerY := float32(state.Board.Height) / float32(2)
	offsetY := float32(head.Y) - centerY
	if offsetY < 0 {
		scorecard.Add(b.UP, m.weight)
	}
	if offsetY > 0 {
		scorecard.Add(b.DOWN, m.weight)
	}
}

// AvoidBiggerSnakes allows a snake to move away from larger snakes
type AvoidBiggerSnakes struct {
	weight Score
}

func (m AvoidBiggerSnakes) move(state b.GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	biggerSnake, err := findNearbySnake(state, head)
	if err == ErrNoBiggerSnakes {
		return // Nothing to do
	}

	// The closer the snake is, the greater the incentive should be to move away
	dist := head.DistanceTo(biggerSnake.Head)
	maxDist := state.Board.Width + state.Board.Height - 2
	weight := m.weight * Score(maxDist-dist)

	// Incentivize moves away from the bigger snake
	if head.X > biggerSnake.Head.X {
		scorecard.Add(b.RIGHT, weight)
	} else {
		scorecard.Add(b.LEFT, weight)
	}
	if head.Y > biggerSnake.Head.Y {
		scorecard.Add(b.UP, weight)
	} else {
		scorecard.Add(b.DOWN, weight)
	}

	debug(state).Msgf("Found bigger snake at %s, %d block(s) away, weight %d", biggerSnake.Head, dist, weight)
}

func findNearbySnake(state b.GameState, head b.Coord) (b.Snake, error) {
	var closestSnake b.Snake
	minDist := math.MaxInt
	for _, snake := range state.Board.Snakes {
		if snake.Length >= state.You.Length && snake.ID != state.You.ID {
			dist := head.DistanceTo(snake.Head)
			if dist < minDist {
				minDist = dist
				closestSnake = snake
			}
		}
	}
	if minDist == math.MaxInt {
		return closestSnake, ErrNoBiggerSnakes
	}
	return closestSnake, nil
}

type AvoidDeadEnds struct {
	weight Score
}

func (a AvoidDeadEnds) move(state b.GameState, scorecard *Scorecard) {
	board := NewBoard(state)
	head := headOfSnake(state)

	spaceRight := availableSpace(head.Right(), board)
	scorecard.Add(b.RIGHT, Score(spaceRight)*a.weight)

	spaceLeft := availableSpace(head.Left(), board)
	scorecard.Add(b.LEFT, Score(spaceLeft)*a.weight)

	spaceDown := availableSpace(head.Down(), board)
	scorecard.Add(b.DOWN, Score(spaceDown)*a.weight)

	spaceUp := availableSpace(head.Up(), board)
	scorecard.Add(b.UP, Score(spaceUp)*a.weight)
}

func availableSpace(start b.Coord, board *Board) int {
	visited := make(map[b.Coord]bool, 0)
	space := 0
	toVisit := make([]b.Coord, 0)
	toVisit = append(toVisit, start)
	for len(toVisit) > 0 {
		curr := toVisit[0]
		toVisit = toVisit[1:]
		// Avoid revisiting squares
		if !visited[curr] {
			visited[curr] = true

			if board.isEmpty(curr) {
				space += 1
				toVisit = append(toVisit, curr.Right(), curr.Left(), curr.Up(), curr.Down())
			}
		}
	}
	return space
}

type Board struct {
	occupied map[b.Coord]bool
	width    int
	height   int
}

func NewBoard(state b.GameState) *Board {
	return &Board{
		occupied: build(state),
		width:    state.Board.Height,
		height:   state.Board.Width,
	}
}

func build(state b.GameState) map[b.Coord]bool {
	board := make(map[b.Coord]bool)

	// Mark all the snakes
	for _, snake := range state.Board.Snakes {
		for _, body := range snake.Body {
			board[body] = true
		}
	}

	// Mark all the hazards
	for _, hazard := range state.Board.Hazards {
		board[hazard] = true
	}

	return board
}

func (b *Board) isEmpty(coord b.Coord) bool {
	if coord.X < 0 {
		return false
	}
	if coord.X >= b.width {
		return false
	}
	if coord.Y < 0 {
		return false
	}
	if coord.Y >= b.height {
		return false
	}
	return !b.occupied[coord]
}

// MoveToFood allows a snake to prefer moves where more food exists.
type MoveToFood struct {
	weight Score
}

func (m MoveToFood) move(state b.GameState, scorecard *Scorecard) {
	var foodToRight, foodToLeft, foodAbove, foodBelow = 0, 0, 0, 0
	head := headOfSnake(state)
	for _, food := range state.Board.Food {
		if food.X > head.X {
			foodToRight += 1
		}
		if food.X < head.X {
			foodToLeft += 1
		}
		if food.Y > head.Y {
			foodAbove += 1
		}
		if food.Y < head.Y {
			foodBelow += 1
		}
	}

	// Update the scorecard
	log.Debug().Msgf("Found food(s) %s=%d, %s=%d, %s=%d, %s=%d",
		b.RIGHT, foodToRight, b.LEFT, foodToLeft, b.UP, foodAbove, b.DOWN, foodBelow)
	scorecard.Add(b.RIGHT, Score(foodToRight)*m.weight)
	scorecard.Add(b.LEFT, Score(foodToLeft)*m.weight)
	scorecard.Add(b.UP, Score(foodAbove)*m.weight)
	scorecard.Add(b.DOWN, Score(foodBelow)*m.weight)
}
