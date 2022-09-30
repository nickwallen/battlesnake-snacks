package snacks

import (
	"errors"
	b "github.com/nickwallen/battlesnake-snacks/battlesnake"
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
	weight float64
}

func (m *MoveToCenter) move(state b.GameState, scorecard *Scorecard) {
	head := headOfSnake(state)
	centerX := float64(state.Board.Width) / float64(2)
	offsetX := float64(head.X) - centerX
	if offsetX < 0 {
		scorecard.Add(b.RIGHT, Score(m.weight*-offsetX))
	}
	if offsetX > 0 {
		scorecard.Add(b.LEFT, Score(m.weight*offsetX))
	}
	centerY := float64(state.Board.Height) / float64(2)
	offsetY := float64(head.Y) - centerY
	if offsetY < 0 {
		scorecard.Add(b.UP, Score(m.weight*-offsetY))
	}
	if offsetY > 0 {
		scorecard.Add(b.DOWN, Score(m.weight*offsetY))
	}
}

// AvoidBiggerSnakes allows a snake to move away from larger snakes
type AvoidBiggerSnakes struct {
	weight float64
}

func (m AvoidBiggerSnakes) move(state b.GameState, scorecard *Scorecard) {
	var rightWeight, leftWeight, aboveWeight, belowWeight = 0.0, 0.0, 0.0, 0.0
	head := headOfSnake(state)
	maxDist := state.Board.Width + state.Board.Height - 2
	for _, snake := range state.Board.Snakes {
		if state.You.Length > snake.Length {
			continue // Ignore smaller snakes
		}
		if snake.ID == state.You.ID {
			continue // Ignore yourself
		}

		// The closer the snake is, the greater the incentive should be to move away
		dist := head.DistanceTo(snake.Head)
		weight := m.weight * float64(maxDist-dist)
		debug(state).Msgf("Found bigger snake at %s, %d block(s) away", snake.Head, dist)

		// Incentivize moves away from the bigger snake
		if head.X > snake.Head.X {
			rightWeight += weight
		} else {
			leftWeight += weight
		}
		if head.Y > snake.Head.Y {
			aboveWeight += weight
		} else {
			belowWeight += weight
		}
	}

	// Update the scorecard
	scorecard.Add(b.RIGHT, Score(rightWeight))
	scorecard.Add(b.LEFT, Score(leftWeight))
	scorecard.Add(b.UP, Score(aboveWeight))
	scorecard.Add(b.DOWN, Score(belowWeight))

	debug(state).Msgf("Moving away from snakes %s=%f, %s=%f, %s=%f, %s=%f",
		b.LEFT, leftWeight, b.RIGHT, rightWeight, b.UP, aboveWeight, b.DOWN, belowWeight)
}

type AvoidDeadEnds struct {
	weight float64
}

func (a AvoidDeadEnds) move(state b.GameState, scorecard *Scorecard) {
	board := NewBoard(state)
	head := headOfSnake(state)

	spaceRight := availableSpace(head.Right(), board) * a.weight
	scorecard.Add(b.RIGHT, Score(spaceRight))

	spaceLeft := availableSpace(head.Left(), board) * a.weight
	scorecard.Add(b.LEFT, Score(spaceLeft))

	spaceBelow := availableSpace(head.Down(), board) * a.weight
	scorecard.Add(b.DOWN, Score(spaceBelow))

	spaceAbove := availableSpace(head.Up(), board) * a.weight
	scorecard.Add(b.UP, Score(spaceAbove))

	debug(state).Msgf("Moving to space %s=%f, %s=%f, %s=%f, %s=%f",
		b.LEFT, spaceLeft, b.RIGHT, spaceRight, b.UP, spaceAbove, b.DOWN, spaceBelow)
}

func availableSpace(start b.Coord, board *Board) float64 {
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
	return float64(space)
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
	weight float64
}

func (m MoveToFood) move(state b.GameState, scorecard *Scorecard) {
	var foodToRight, foodToLeft, foodAbove, foodBelow = 0.0, 0.0, 0.0, 0.0
	head := headOfSnake(state)
	maxDist := state.Board.Width + state.Board.Height - 2
	for _, food := range state.Board.Food {
		// The closer the food, the greater the weight
		foodWeight := float64(maxDist-head.DistanceTo(food)) * m.weight
		if food.X > head.X {
			foodToRight += foodWeight
		}
		if food.X < head.X {
			foodToLeft += foodWeight
		}
		if food.Y > head.Y {
			foodAbove += foodWeight
		}
		if food.Y < head.Y {
			foodBelow += foodWeight
		}
	}

	// Update the scorecard
	scorecard.Add(b.RIGHT, Score(foodToRight))
	scorecard.Add(b.LEFT, Score(foodToLeft))
	scorecard.Add(b.UP, Score(foodAbove))
	scorecard.Add(b.DOWN, Score(foodBelow))

	debug(state).Msgf("Moving to food %s=%f, %s=%f, %s=%f, %s=%f",
		b.LEFT, foodToLeft, b.RIGHT, foodToRight, b.UP, foodAbove, b.DOWN, foodBelow)
}
