package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Right(t *testing.T) {
	start := Coord{-1, 2}
	require.EqualValues(t, Coord{0, 2}, start.Right())
	require.EqualValues(t, Coord{1, 2}, start.Right().Right())
	require.EqualValues(t, Coord{2, 2}, start.Right().Right().Right())
	require.EqualValues(t, Coord{3, 2}, start.Right().Right().Right().Right())
	require.EqualValues(t, Coord{4, 2}, start.Right().Right().Right().Right().Right())
}

func Test_Left(t *testing.T) {
	start := Coord{4, 2}
	require.EqualValues(t, Coord{3, 2}, start.Left())
	require.EqualValues(t, Coord{2, 2}, start.Left().Left())
	require.EqualValues(t, Coord{1, 2}, start.Left().Left().Left())
	require.EqualValues(t, Coord{0, 2}, start.Left().Left().Left().Left())
	require.EqualValues(t, Coord{-1, 2}, start.Left().Left().Left().Left().Left())
}

func Test_Up(t *testing.T) {
	start := Coord{2, -1}
	require.EqualValues(t, Coord{2, 0}, start.Up())
	require.EqualValues(t, Coord{2, 1}, start.Up().Up())
	require.EqualValues(t, Coord{2, 2}, start.Up().Up().Up())
	require.EqualValues(t, Coord{2, 3}, start.Up().Up().Up().Up())
	require.EqualValues(t, Coord{2, 4}, start.Up().Up().Up().Up().Up())
}

func Test_Down(t *testing.T) {
	start := Coord{2, 4}
	require.EqualValues(t, Coord{2, 3}, start.Down())
	require.EqualValues(t, Coord{2, 2}, start.Down().Down())
	require.EqualValues(t, Coord{2, 1}, start.Down().Down().Down())
	require.EqualValues(t, Coord{2, 0}, start.Down().Down().Down().Down())
	require.EqualValues(t, Coord{2, -1}, start.Down().Down().Down().Down().Down())
}

func Test_DistanceTo(t *testing.T) {
	start := Coord{-2, -2}
	end := Coord{2, 2}
	require.Equal(t, 8, start.DistanceTo(end))
	require.Equal(t, 8, end.DistanceTo(start))
}

func Test_MoveTo(t *testing.T) {
	start := Coord{0, 0}

	require.Equal(t, RIGHT, start.MoveTo(Coord{1, 0}))
	require.Equal(t, UP, start.MoveTo(Coord{1, 1}))
	require.Equal(t, UP, start.MoveTo(Coord{0, 1}))
	require.Equal(t, UP, start.MoveTo(Coord{-1, 1}))
	require.Equal(t, LEFT, start.MoveTo(Coord{-1, 0}))
	require.Equal(t, DOWN, start.MoveTo(Coord{-1, -1}))
	require.Equal(t, DOWN, start.MoveTo(Coord{0, -1}))
	require.Equal(t, DOWN, start.MoveTo(Coord{1, -1}))
}
