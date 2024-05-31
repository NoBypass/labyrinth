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

var (
	gridSize int
	size     int
)

func init() {
	s := flag.Int("size", 15, "size of the labyrinth")
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

	fmt.Println("Grid generated in", time.Since(start).String())

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	grid[0][0].borders[Left].shown = false
	grid[gridSize-1][gridSize-1].borders[Right].shown = false

	beenTo := make(map[*field]struct{}, gridSize*gridSize)
	fieldsWithOptions := make([]*field, 0, gridSize*gridSize)
	var backtracks int
	var loops int
	pointer := grid[0][0]

	for len(beenTo) < gridSize*gridSize {
		loops++
		available := make(map[direction]*field, 4)
		availableDirs := make([]direction, 0, len(available))

		for i := range direction(4) {
			potentialPointer := pointer.move(grid, i)
			if potentialPointer != nil {
				if _, hasBeenTo := beenTo[potentialPointer]; !hasBeenTo {
					available[i] = potentialPointer
					availableDirs = append(availableDirs, i)
				}
			}
		}

		beenTo[pointer] = struct{}{}
		if len(availableDirs) == 0 {
			backtracks++
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
	fmt.Println("Labyrinth generated in", time.Since(start).String())
	fmt.Println("Backtracks:", backtracks)
	fmt.Println("Loops:", loops)
}
