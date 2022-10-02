package battlesnake

import (
	"fmt"
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c Coord) Move(move Move) Coord {
	switch move {
	case RIGHT:
		return c.Right()
	case LEFT:
		return c.Left()
	case UP:
		return c.Up()
	case DOWN:
		return c.Down()
	default:
		panic("unexpected move")
	}
}

func (c Coord) Right() Coord {
	return Coord{
		X: c.X + 1,
		Y: c.Y,
	}
}

func (c Coord) Left() Coord {
	return Coord{
		X: c.X - 1,
		Y: c.Y,
	}
}

func (c Coord) Up() Coord {
	return Coord{
		X: c.X,
		Y: c.Y + 1,
	}
}

func (c Coord) Down() Coord {
	return Coord{
		X: c.X,
		Y: c.Y - 1,
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// DistanceTo Returns the Manhattan Distance to the target
func (c Coord) DistanceTo(target Coord) int {
	X := abs(c.X - target.X)
	Y := abs(c.Y - target.Y)
	return X + Y
}

func (c Coord) MoveTo(target Coord) Move {
	xDelta := abs(target.X - c.X)
	yDelta := abs(target.Y - c.Y)
	if xDelta > yDelta {
		if c.X < target.X {
			return RIGHT
		}
		return LEFT
	} else {
		if c.Y < target.Y {
			return UP
		}
		return DOWN
	}
}

func (c Coord) MoveAway(target Coord) Move {
	xDelta := abs(target.X - c.X)
	yDelta := abs(target.Y - c.Y)
	if xDelta > yDelta {
		if c.X < target.X {
			return LEFT
		}
		return RIGHT
	} else {
		if c.Y < target.Y {
			return DOWN
		}
		return UP
	}
}
