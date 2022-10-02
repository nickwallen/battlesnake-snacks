package snacks

import (
	"errors"
	b "github.com/nickwallen/battlesnake-snacks/internal/battlesnake"
	"math"
)

var (
	ErrNoBiggerSnakes = errors.New("no bigger snakes found")
	ErrNoFood         = errors.New("no food found")
)

// StayInBounds A strategy that keeps the snake within the game boundaries.
type StayInBounds struct {
}

func (s *StayInBounds) move(state b.GameState, card *Scorecard) {
	scorecard := NewLoggingScorecard("stay-in-bounds", state, card)
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

func (a *NoCollisions) move(state b.GameState, card *Scorecard) {
	var avoid []b.Coord

	// Avoid self collisions, other snakes and hazards
	for _, square := range state.You.Body {
		avoid = append(avoid, square)
	}
	for _, opponent := range state.Board.Snakes {
		for _, square := range opponent.Body {
			avoid = append(avoid, square)
		}
	}
	for _, hazard := range state.Board.Hazards {
		avoid = append(avoid, hazard)
	}

	scorecard := NewLoggingScorecard("no-collisions", state, card)
	head := headOfSnake(state)
	for _, square := range avoid {
		if head.Right() == square {
			scorecard.Unsafe(b.RIGHT)
		}
		if head.Left() == square {
			scorecard.Unsafe(b.LEFT)
		}
		if head.Up() == square {
			scorecard.Unsafe(b.UP)
		}
		if head.Down() == square {
			scorecard.Unsafe(b.DOWN)
		}
	}
}

// MoveToClosestFood allows a snake to move toward the nearest food source.
type MoveToClosestFood struct {
	weight Score
}

func (m MoveToClosestFood) move(state b.GameState, card *Scorecard) {
	head := headOfSnake(state)
	closestFood, err := findNearbyFood(state, head)
	if err == ErrNoFood {
		return // Nothing to do
	}

	// Incentivize moves that take us closer to the food
	scorecard := NewLoggingScorecard("move-to-closest-food", state, card)
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

func (m *MoveToCenter) move(state b.GameState, card *Scorecard) {
	scorecard := NewLoggingScorecard("move-to-center", state, card)
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

// MoveToWalls allows a snake to remain close to the outer walls.
type MoveToWalls struct {
	weight float64
}

func (m *MoveToWalls) move(state b.GameState, card *Scorecard) {
	scorecard := NewLoggingScorecard("move-to-walls", state, card)
	head := headOfSnake(state)
	centerX := float64(state.Board.Width) / float64(2)
	offsetX := float64(head.X) - centerX
	if offsetX < 0 {
		scorecard.Add(b.LEFT, Score(m.weight*-offsetX))
	}
	if offsetX > 0 {
		scorecard.Add(b.RIGHT, Score(m.weight*offsetX))
	}
	centerY := float64(state.Board.Height) / float64(2)
	offsetY := float64(head.Y) - centerY
	if offsetY < 0 {
		scorecard.Add(b.DOWN, Score(m.weight*-offsetY))
	}
	if offsetY > 0 {
		scorecard.Add(b.UP, Score(m.weight*offsetY))
	}
}

// AvoidBiggerSnakes allows a snake to move away from larger snakes
type AvoidBiggerSnakes struct {
	weight float64
}

func (m AvoidBiggerSnakes) move(state b.GameState, card *Scorecard) {
	var weightRight, weightLeft, weightUp, weightDown = 0.0, 0.0, 0.0, 0.0
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
			weightRight += weight
		} else {
			weightLeft += weight
		}
		if head.Y > snake.Head.Y {
			weightUp += weight
		} else {
			weightDown += weight
		}
	}

	// Update the scorecard
	scorecard := NewLoggingScorecard("avoid-bigger-snakes", state, card)
	scorecard.Add(b.RIGHT, Score(weightRight))
	scorecard.Add(b.LEFT, Score(weightLeft))
	scorecard.Add(b.UP, Score(weightUp))
	scorecard.Add(b.DOWN, Score(weightDown))
}

// AvoidDeadEnds allows a snake to avoid dead ends.
type AvoidDeadEnds struct {
}

func (m AvoidDeadEnds) move(state b.GameState, card *Scorecard) {
	board := NewBoard(state)
	head := headOfSnake(state)
	scorecard := NewLoggingScorecard("avoid-dead-ends", state, card)
	spaceLeft := availableSpace(head.Left(), board)
	if spaceLeft < state.You.Length {
		scorecard.Unsafe(b.LEFT)
		debug(state).Msgf("Dead-end left! Have %d square(s), need %d", spaceLeft, state.You.Length)
	}

	spaceRight := availableSpace(head.Right(), board)
	if spaceRight < state.You.Length {
		scorecard.Unsafe(b.RIGHT)
		debug(state).Msgf("Dead-end right! Have %d square(s), need %d", spaceRight, state.You.Length)
	}

	spaceUp := availableSpace(head.Up(), board)
	if spaceUp < state.You.Length {
		scorecard.Unsafe(b.UP)
		debug(state).Msgf("Dead-end up! Have %d square(s), need %d", spaceRight, state.You.Length)
	}

	spaceDown := availableSpace(head.Down(), board)
	if spaceDown < state.You.Length {
		scorecard.Unsafe(b.DOWN)
		debug(state).Msgf("Dead-end down! Have %d square(s), need %d", spaceRight, state.You.Length)
	}
}

// MoveToSpace allows a snake to move towards areas with more available space.
type MoveToSpace struct {
	weight float64
}

func (a MoveToSpace) move(state b.GameState, card *Scorecard) {
	board := NewBoard(state)
	head := headOfSnake(state)
	totalSpaces := state.Board.Height * state.Board.Width
	scorecard := NewLoggingScorecard("move-to-space", state, card)

	weightRight := float64(availableSpace(head.Right(), board)) / float64(totalSpaces) * 10 * a.weight
	scorecard.Add(b.RIGHT, Score(weightRight))

	weightLeft := float64(availableSpace(head.Left(), board)) / float64(totalSpaces) * 10 * a.weight
	scorecard.Add(b.LEFT, Score(weightLeft))

	weightDown := float64(availableSpace(head.Down(), board)) / float64(totalSpaces) * 10 * a.weight
	scorecard.Add(b.DOWN, Score(weightDown))

	weightUp := float64(availableSpace(head.Up(), board)) / float64(totalSpaces) * 10 * a.weight
	scorecard.Add(b.UP, Score(weightUp))
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
	weight float64
}

func (m MoveToFood) move(state b.GameState, card *Scorecard) {
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
	scorecard := NewLoggingScorecard("move-to-food", state, card)
	scorecard.Add(b.RIGHT, Score(foodToRight))
	scorecard.Add(b.LEFT, Score(foodToLeft))
	scorecard.Add(b.UP, Score(foodAbove))
	scorecard.Add(b.DOWN, Score(foodBelow))
}

type AttackSmallerSnakes struct {
	weight float64
}

func (a AttackSmallerSnakes) move(state b.GameState, card *Scorecard) {
	var weightRight, weightLeft, weightUp, weightDown = 0.0, 0.0, 0.0, 0.0
	head := headOfSnake(state)
	maxDist := state.Board.Width + state.Board.Height - 2
	for _, snake := range state.Board.Snakes {
		if state.You.Length <= snake.Length {
			continue // Ignore bigger snakes
		}
		if snake.ID == state.You.ID {
			continue // Ignore yourself
		}

		// The closer the snake is, the greater the incentive should be to attack
		dist := head.DistanceTo(snake.Head)
		weight := a.weight * float64(maxDist-dist)
		debug(state).Msgf("Found smaller snake at %s, %d block(s) away", snake.Head, dist)

		// Incentivize moves toward the smaller snake
		if head.X > snake.Head.X {
			weightLeft += weight
		} else {
			weightRight += weight
		}
		if head.Y > snake.Head.Y {
			weightDown += weight
		} else {
			weightUp += weight
		}
	}

	// Update the scorecard
	scorecard := NewLoggingScorecard("attack-smaller-snakes", state, card)
	scorecard.Add(b.RIGHT, Score(weightRight))
	scorecard.Add(b.LEFT, Score(weightLeft))
	scorecard.Add(b.UP, Score(weightUp))
	scorecard.Add(b.DOWN, Score(weightDown))
}
