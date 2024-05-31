package main

import (
	"flag"
	"fmt"
	"github.com/fogleman/gg"
	"math/rand"
	"time"
)

const multiplier = 25
const offset = 50

var moves = []move{
	{-1, 0},
	{1, 0},
	{0, -1},
	{0, 1},
}

var (
	gridSize int
	size     int
)

func init() {
	s := flag.Int("size", 10, "size of the labyrinth")
	flag.Parse()
	gridSize = *s
	size = gridSize * multiplier
}

func main() {
	start := time.Now()
	dc := gg.NewContext(size+offset*2, size+offset*2)

	grid := make([][]*field, gridSize)
	lines := make(map[*line]struct{}, (gridSize+1)*gridSize*2)

	for y := range grid {
		grid[y] = make([]*field, gridSize)
		for x := range grid[y] {
			f := &field{
				x: x,
				y: y,
			}
			f.borders = map[direction]*line{
				Right: newLine(Right, f),
				Left:  newLine(Left, f),
				Down:  newLine(Down, f),
				Up:    newLine(Up, f),
			}
			if y != 0 {
				f.borders[Up] = grid[y-1][x].borders[Down]
			}
			if x != 0 {
				f.borders[Left] = grid[y][x-1].borders[Right]
			}
			for _, l := range f.borders {
				lines[l] = struct{}{}
			}
			grid[y][x] = f
		}
	}

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	grid[0][0].borders[Left].shown = false
	grid[gridSize-1][gridSize-1].borders[Right].shown = false

	beenTo := make(map[*field]struct{}, gridSize*gridSize)
	fieldsWithOptions := make([]*field, 0, gridSize*gridSize)
	pointer := grid[0][0]

	for len(beenTo) < gridSize*gridSize {
		available := make(map[direction]*field, 4)
		availableDirs := make([]direction, 0, len(available))

		for i, m := range moves {
			newX, newY := pointer.x+m.dx, pointer.y+m.dy
			if newX >= 0 && newX < gridSize && newY >= 0 && newY < gridSize {
				if _, hasBeenTo := beenTo[grid[newY][newX]]; hasBeenTo {
					continue
				}
				available[direction(i)] = grid[newY][newX]
				availableDirs = append(availableDirs, direction(i))
			}
		}

		beenTo[pointer] = struct{}{}
		if len(availableDirs) == 0 {
			pointer = fieldsWithOptions[0]
			fieldsWithOptions = fieldsWithOptions[1:]
			continue
		} else if len(availableDirs) > 1 {
			fieldsWithOptions = append(fieldsWithOptions, pointer)
		}
		dir := availableDirs[rand.Intn(len(availableDirs))]

		pointer.borders[dir].shown = false
		pointer = available[dir]
	}

	for l := range lines {
		l.draw(dc)
	}

	dc.SavePNG("labyrinth.png")
	fmt.Printf("\nTook %s\n", time.Since(start).String())
}
