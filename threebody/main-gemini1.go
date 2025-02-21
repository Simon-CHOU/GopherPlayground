package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	trailLength  = 30 * 60 // 30 seconds in frames
)

type Body struct {
	x, y      float64
	vx, vy    float64
	mass      float64
	radius    float64
	trail     []vector.Point
	trailLife []int
	collided  bool
}

var bodies [3]Body

func init() {
	rand.Seed(time.Now().UnixNano())
	initBodies()
}

func initBodies() {
	for i := range bodies {
		for {
			bodies[i].x = rand.Float64() * screenWidth
			bodies[i].y = rand.Float64() * screenHeight
			bodies[i].vx = (rand.Float64() - 0.5) * 10
			bodies[i].vy = (rand.Float64() - 0.5) * 10
			bodies[i].mass = rand.Float64()*5 + 1            // Mass between 1 and 6
			bodies[i].radius = math.Sqrt(bodies[i].mass) * 5 // Radius scales with mass

			// Ensure bodies are not too close to each other
			minDistance := 50.0
			validPosition := true
			for j := 0; j < i; j++ {
				dx := bodies[i].x - bodies[j].x
				dy := bodies[i].y - bodies[j].y
				distance := math.Sqrt(dx*dx + dy*dy)
				if distance < minDistance {
					validPosition = false
					break
				}
			}
			if validPosition {
				break
			}
		}
		bodies[i].trail = make([]vector.Point, 0, trailLength)
		bodies[i].trailLife = make([]int, 0, trailLength)
	}
}

func (b *Body) update(bodies []Body) {
	for i := range bodies {
		if &bodies[i] == b {
			continue
		}
		dx := bodies[i].x - b.x
		dy := bodies[i].y - b.y
		distance := math.Sqrt(dx*dx + dy*dy)
		force := bodies[i].mass * b.mass / (distance * distance)
		b.vx += force * dx / distance / b.mass
		b.vy += force * dy / distance / b.mass

		// Collision detection
		if distance < b.radius+bodies[i].radius && !b.collided && !bodies[i].collided {
			b.collided = true
			bodies[i].collided = true
			initBodies() // Restart simulation on collision
			return
		}
	}

	b.x += b.vx
	b.y += b.vy

	// Wrap around edges
	if b.x < 0 {
		b.x = screenWidth
	} else if b.x > screenWidth {
		b.x = 0
	}
	if b.y < 0 {
		b.y = screenHeight
	} else if b.y > screenHeight {
		b.y = 0
	}

	// Update trail
	b.trail = append(b.trail, vector.Point{float32(b.x), float32(b.y)})
	b.trailLife = append(b.trailLife, trailLength)
	if len(b.trail) > trailLength {
		b.trail = b.trail[1:]
		b.trailLife = b.trailLife[1:]
	}

	// Fade trail
	for i := range b.trailLife {
		b.trailLife[i]--
	}
}

type Game struct{}

func (g *Game) Update() error {
	for i := range bodies {
		bodies[i].update(bodies[:])
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range bodies {
		// Draw trail
		for i, p := range b.trail {
			alpha := float32(b.trailLife[i]) / trailLength
			c := color.RGBA{255, 255, 255, uint8(alpha * 255)} // White trail fading out
			vector.FillCircle(screen, p, 2, c)
		}
		vector.FillCircle(screen, vector.Point{float32(b.x), float32(b.y)}, float32(b.radius), color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three-Body Problem")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
