package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	G            = 6.67430e-11 // 引力常数
	dt           = 0.1         // 时间步长
)

type Body struct {
	x, y, vx, vy, mass float64
}

type Game struct {
	bodies [3]Body
}

func (g *Game) Update() error {
	for i := range g.bodies {
		g.updateBody(i)
	}
	if g.checkOutOfBounds() {
		g.resetBodies()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, body := range g.bodies {
		vector.DrawFilledCircle(screen, float32(body.x), float32(body.y), 10, color.White, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) updateBody(i int) {
	body := &g.bodies[i]
	ax, ay := 0.0, 0.0

	for j, other := range g.bodies {
		if i == j {
			continue
		}
		dx := other.x - body.x
		dy := other.y - body.y
		r := math.Sqrt(dx*dx + dy*dy)
		f := G * body.mass * other.mass / (r * r)
		ax += f * dx / r / body.mass
		ay += f * dy / r / body.mass
	}

	body.vx += ax * dt
	body.vy += ay * dt
	body.x += body.vx * dt
	body.y += body.vy * dt
}

func (g *Game) checkOutOfBounds() bool {
	for _, body := range g.bodies {
		if body.x < 0 || body.x > screenWidth || body.y < 0 || body.y > screenHeight {
			return true
		}
	}
	return false
}

func (g *Game) resetBodies() {
	for i := range g.bodies {
		g.bodies[i] = Body{
			x:    rand.Float64() * screenWidth,
			y:    rand.Float64() * screenHeight,
			vx:   (rand.Float64() - 0.5) * 100,
			vy:   (rand.Float64() - 0.5) * 100,
			mass: rand.Float64()*1e12 + 1e11,
		}
	}
}

func main() {
	game := &Game{}
	game.resetBodies()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three-Body Problem Simulation")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
