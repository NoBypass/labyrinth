package main

import "github.com/fogleman/gg"

type direction int
type rgb struct {
	r, g, b float64
}

type move struct {
	dy, dx int
}

const (
	Up direction = iota
	Down
	Left
	Right
)

var moves = []move{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

type line struct {
	color          rgb
	shown          bool
	x1, y1, x2, y2 int
}

func (f *field) move(grid [][]*field, d direction) *field {
	m := moves[d]
	newX, newY := f.x+m.dx, f.y+m.dy
	if newX >= 0 && newX < gridSize && newY >= 0 && newY < gridSize {
		return grid[newY][newX]
	}
	return nil
}

type field struct {
	x, y    int
	borders map[direction]*line
}

func (l *line) draw(dc *gg.Context) {
	if l.shown {
		dc.SetRGB(l.color.r, l.color.g, l.color.b)
		dc.DrawLine(float64(l.x1), float64(l.y1), float64(l.x2), float64(l.y2))
		dc.Stroke()
	}
}

func newLine(d direction, f *field, colors ...rgb) *line {
	c := rgb{0, 0, 0}
	if len(colors) > 0 {
		c = colors[0]
	}

	l := &line{
		shown: true,
		color: c,
	}
	switch d {
	case Up:
		l.x1 = f.x
		l.y1 = f.y - 1
		l.x2 = f.x + 1
		l.y2 = f.y - 1
	case Down:
		l.x1 = f.x
		l.y1 = f.y
		l.x2 = f.x + 1
		l.y2 = f.y
	case Left:
		l.x1 = f.x
		l.y1 = f.y
		l.x2 = f.x
		l.y2 = f.y - 1
	case Right:
		l.x1 = f.x + 1
		l.y1 = f.y
		l.x2 = f.x + 1
		l.y2 = f.y - 1
	}

	l.y1 += 1
	l.y2 += 1

	l.x1 *= multiplier
	l.y1 *= multiplier
	l.x2 *= multiplier
	l.y2 *= multiplier

	l.x1 += offset
	l.y1 += offset
	l.x2 += offset
	l.y2 += offset
	return l
}
