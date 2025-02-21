package main

import (
	"image/color"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	trailLength  = 30 * 60 // 30 seconds * 60 frames per second
	bodyRadius   = 10
)

type Body struct {
	x, y, vx, vy, mass float64
	color              color.RGBA
	trail              []*TrailPoint
}

type TrailPoint struct {
	x, y float64
	age  int
}

func (b *Body) attract(other *Body) (fx, fy float64) {
	dx := other.x - b.x
	dy := other.y - b.y
	dist := math.Sqrt(dx*dx + dy*dy)
	force := 0.1 * b.mass * other.mass / (dist*dist + 0.1) // 避免除以0，并调整引力系数
	fx = force * dx / dist
	fy = force * dy / dist
	return
}

func (b *Body) update(bodies []*Body) {
	fx, fy := 0.0, 0.0
	for _, other := range bodies {
		if other != b {
			fxx, fyy := b.attract(other)
			fx += fxx
			fy += fyy
		}
	}

	b.vx += fx / b.mass
	b.vy += fy / b.mass

	b.x += b.vx
	b.y += b.vy

	// 添加尾迹点
	b.trail = append(b.trail, &TrailPoint{x: b.x, y: b.y, age: 0})
	if len(b.trail) > trailLength {
		b.trail = b.trail[1:]
	}

	// 更新尾迹年龄
	for _, tp := range b.trail {
		tp.age++
	}
}

type Game struct {
	bodies []*Body
}

func (g *Game) Update() error {
	for _, b := range g.bodies {
		b.update(g.bodies)
	}

	// 碰撞检测
	for i := 0; i < len(g.bodies); i++ {
		for j := i + 1; j < len(g.bodies); j++ {
			dx := g.bodies[i].x - g.bodies[j].x
			dy := g.bodies[i].y - g.bodies[j].y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < bodyRadius*2 { // 碰撞半径的两倍
				g.initBodies() // 重新开始模拟
				return nil
			}
		}
	}

	// 边界检测
	for _, b := range g.bodies {
		if b.x < 0 || b.x > screenWidth || b.y < 0 || b.y > screenHeight {
			g.initBodies() // 重新开始模拟
			return nil
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	for _, b := range g.bodies {
		// 绘制尾迹
		for _, tp := range b.trail {
			opacity := float32(trailLength-tp.age) / trailLength
			c := b.color
			c.A = uint8(255 * opacity)
			ebitenutil.DrawRect(screen, float32(tp.x)-1, float32(tp.y)-1, 2, 2, c) // 更细的尾迹
		}

		// 绘制球体
		ebitenutil.DrawCircle(screen, float32(b.x), float32(b.y), bodyRadius, b.color) // 使用 bodyRadius
	}
}

func (g *Game) initBodies() {
	rand.Seed(time.Now().UnixNano())
	g.bodies = []*Body{
		{x: float64(screenWidth) * rand.Float64(), y: float64(screenHeight) * rand.Float64(), vx: rand.Float64()*2 - 1, vy: rand.Float64()*2 - 1, mass: 1, color: color.RGBA{255, 0, 0, 255}, trail: []*TrailPoint{}},
		{x: float64(screenWidth) * rand.Float64(), y: float64(screenHeight) * rand.Float64(), vx: rand.Float64()*2 - 1, vy: rand.Float64()*2 - 1, mass: 1, color: color.RGBA{0, 255, 0, 255}, trail: []*TrailPoint{}},
		{x: float64(screenWidth) * rand.Float64(), y: float64(screenHeight) * rand.Float64(), vx: rand.Float64()*2 - 1, vy: rand.Float64()*2 - 1, mass: 1, color: color.RGBA{0, 0, 255, 255}, trail: []*TrailPoint{}},
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{}
	game.initBodies()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three-Body Problem")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
