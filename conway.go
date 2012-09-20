package main

import (
	"conway/universe"
	"fmt"
	"math/rand"
	"time"
)

const (
	rows          = 50
	cols          = 50
	generations   = 1000
	threshold     = 50
	dead          = 0
	alive         = 1
	erase_display = "\x1B[2J"
	fmt_inverse   = "\x1B[7m"
	fmt_normal    = "\x1B[0m"
)

type Universe struct {
	cols, rows int
	content    []int
}

// Gets the cell at position i, j
func (universe *Universe) get(i, j int) int {
	return universe.content[i*universe.cols+j]
}

// Sets the cell at position i, j to value
func (universe *Universe) set(i, j, value int) {
	universe.content[i*universe.cols+j] = value
}

// SumNeighbours counts the amount of live cells surrounding i, j
func (universe *Universe) sumNeighbours(i, j int) int {
	sum := 0

	cols, rows := universe.cols, universe.rows

	for dx := (i - 1 + rows) % rows; dx <= (i+1+rows)%rows; dx++ {
		for dy := (j - 1 + cols) % cols; dy <= (j+1+cols)%cols; dy++ {
			if dx != i || dy != j {
				sum += universe.get(dx, dy)
			}
		}
	}
	return sum
}

// Lives returns whether the cell should live or die based on the result
// from sumNeighbours
func (universe *Universe) Lives(i, j int) bool {
	neighbours := universe.sumNeighbours(i, j)
	return neighbours == 3 || (universe.get(i, j) == 1 && neighbours == 2)
}

// NextGeneration returns a new Universe with the state of the next iteration
// of the universe
func (universe *Universe) NextGeneration() *Universe {
	newUniverse := New(universe.rows, universe.cols)

	for i := 0; i < universe.rows; i++ {
		for j := 0; j < universe.cols; j++ {
			if universe.Lives(i, j) {
				newUniverse.set(i, j, alive)
			} else {
				newUniverse.set(i, j, dead)
			}
		}
	}
	return newUniverse
}

// Shows the current universe using the full terminal screen
func (universe *Universe) Show() {
	// ANSI escape code to clear the terminal screen
	fmt.Print(erase_display)

	// Print cells using ANSI colors
	for i := 0; i < universe.rows; i++ {
		for j := 0; j < universe.cols; j++ {
			if universe.get(i, j) == alive {
				fmt.Print(fmt_inverse + "  " + fmt_normal)
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Print("\n")
	}
}

// Populate fills a board with dead and live cells, based on a set threshold
func (universe *Universe) Populate(threshold int) {
	for i := 0; i < universe.rows; i++ {
		for j := 0; j < universe.cols; j++ {
			if value := rand.Intn(100); value > threshold {
				universe.set(i, j, alive)
			} else {
				universe.set(i, j, dead)
			}
		}
	}
}

// New returns a new Conway's Game of Life universe
func New(rows, cols int) *Universe {
	universe := &Universe{cols: cols, rows: rows, content: make([]int, cols*rows)}
	return universe
}

var ch = make(chan int, 1)

// Program entry point
func main() {
	rand.Seed(time.Now().Unix())

	u := universe.New(rows, cols)

	u.Populate(threshold)
	u.Show()

	for i := 0; i < generations; i++ {
		go func() {
			u = u.NextGeneration()
			ch <- 1
		}()

		time.Sleep(100 * time.Millisecond)
		<-ch

		u.Show()
	}
}
