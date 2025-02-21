package main

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	width    = 100
	height   = 100
	cellSize = 8
)

type Grid [height][width]bool

func main() {
	myApp := app.New()
	window := myApp.NewWindow("生命游戏")

	grid := newGrid()
	content := container.New(NewGridLayout(width))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rect := canvas.NewRectangle(color.Black)
			rect.Resize(fyne.NewSize(cellSize, cellSize))
			content.Add(rect)
		}
	}

	window.SetContent(content)
	window.Resize(fyne.NewSize(width*cellSize, height*cellSize))

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			grid = nextGeneration(grid)
			updateGUI(content, grid)
		}
	}()

	window.ShowAndRun()
}

func newGrid() Grid {
	var grid Grid
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			grid[y][x] = rand.Float32() < 0.3
		}
	}
	return grid
}

func nextGeneration(grid Grid) Grid {
	var newGrid Grid
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			neighbors := countNeighbors(grid, x, y)
			if grid[y][x] {
				newGrid[y][x] = neighbors == 2 || neighbors == 3
			} else {
				newGrid[y][x] = neighbors == 3
			}
		}
	}
	return newGrid
}

func countNeighbors(grid Grid, x, y int) int {
	count := 0
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			nx, ny := (x+dx+width)%width, (y+dy+height)%height
			if grid[ny][nx] {
				count++
			}
		}
	}
	return count
}

func updateGUI(content *fyne.Container, grid Grid) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			rect := content.Objects[y*width+x].(*canvas.Rectangle)
			if grid[y][x] {
				rect.FillColor = color.White
			} else {
				rect.FillColor = color.Black
			}
			rect.Refresh()
		}
	}
}
