package main

import (
	"image/color"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
	G            = 0.1 // 引力常数（为了模拟效果已缩放）
	dt           = 1.0 / 120.0
	bodyRadius   = 10.0
	trailLife    = 30.0 // 轨迹存活时间（秒）
)

// Body 表示一个天体
type Body struct {
	mass     float64
	position Vec2
	velocity Vec2
}

// Vec2 表示二维向量
type Vec2 struct {
	x, y float64
}

// Trail 表示运动轨迹点
type Trail struct {
	position Vec2
	created  time.Time
}

// Game 表示游戏状态
type Game struct {
	bodies []Body
	trails [][]Trail
}

func (v Vec2) add(other Vec2) Vec2 {
	return Vec2{v.x + other.x, v.y + other.y}
}

func (v Vec2) sub(other Vec2) Vec2 {
	return Vec2{v.x - other.x, v.y - other.y}
}

func (v Vec2) mult(s float64) Vec2 {
	return Vec2{v.x * s, v.y * s}
}

func (v Vec2) length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y)
}

func (v Vec2) normalize() Vec2 {
	l := v.length()
	if l == 0 {
		return Vec2{0, 0}
	}
	return Vec2{v.x / l, v.y / l}
}

func NewGame() *Game {
	// 初始化三个天体
	bodies := []Body{
		{
			mass:     1.0,
			position: Vec2{screenWidth * 0.3, screenHeight * 0.3},
			velocity: Vec2{50, 20},
		},
		{
			mass:     1.0,
			position: Vec2{screenWidth * 0.7, screenHeight * 0.3},
			velocity: Vec2{-30, 40},
		},
		{
			mass:     1.0,
			position: Vec2{screenWidth * 0.5, screenHeight * 0.7},
			velocity: Vec2{-20, -30},
		},
	}

	return &Game{
		bodies: bodies,
		trails: make([][]Trail, 3),
	}
}

func (g *Game) Reset() {
	*g = *NewGame()
}

func (g *Game) Update() error {
	// 更新天体位置和速度
	newBodies := make([]Body, len(g.bodies))
	copy(newBodies, g.bodies)

	// 计算加速度并更新速度和位置
	for i := range g.bodies {
		acceleration := Vec2{0, 0}
		for j := range g.bodies {
			if i != j {
				r := g.bodies[j].position.sub(g.bodies[i].position)
				dist := r.length()

				// 检查碰撞
				if dist < bodyRadius*2 {
					g.Reset()
					return nil
				}

				// 计算引力
				force := G * g.bodies[i].mass * g.bodies[j].mass / (dist * dist)
				acceleration = acceleration.add(r.normalize().mult(force / g.bodies[i].mass))
			}
		}

		// 更新速度和位置
		newBodies[i].velocity = g.bodies[i].velocity.add(acceleration.mult(dt))
		newBodies[i].position = g.bodies[i].position.add(newBodies[i].velocity.mult(dt))

		// 检查是否超出边界
		if newBodies[i].position.x < 0 || newBodies[i].position.x > screenWidth ||
			newBodies[i].position.y < 0 || newBodies[i].position.y > screenHeight {
			g.Reset()
			return nil
		}
	}

	// 更新轨迹
	now := time.Now()
	for i := range g.bodies {
		g.trails[i] = append(g.trails[i], Trail{
			position: g.bodies[i].position,
			created:  now,
		})

		// 移除过期的轨迹点
		for len(g.trails[i]) > 0 && now.Sub(g.trails[i][0].created).Seconds() > trailLife {
			g.trails[i] = g.trails[i][1:]
		}
	}

	g.bodies = newBodies
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 绘制背景
	screen.Fill(color.RGBA{0, 0, 0, 255})

	// 绘制轨迹
	now := time.Now()
	for i, bodyTrail := range g.trails {
		baseColor := []color.RGBA{
			{255, 0, 0, 255},
			{0, 255, 0, 255},
			{0, 0, 255, 255},
		}[i]

		for _, trail := range bodyTrail {
			age := now.Sub(trail.created).Seconds()
			alpha := uint8(255 * (1 - age/trailLife))
			trailColor := color.RGBA{
				baseColor.R,
				baseColor.G,
				baseColor.B,
				alpha,
			}

			x, y := trail.position.x, trail.position.y
			screen.Set(int(x), int(y), trailColor)
		}
	}

	// 绘制天体
	colors := []color.RGBA{
		{255, 100, 100, 255},
		{100, 255, 100, 255},
		{100, 100, 255, 255},
	}

	for i, body := range g.bodies {
		// 绘制实心圆
		for dy := -bodyRadius; dy <= bodyRadius; dy++ {
			for dx := -bodyRadius; dx <= bodyRadius; dx++ {
				if dx*dx+dy*dy <= bodyRadius*bodyRadius {
					x := int(body.position.x + dx)
					y := int(body.position.y + dy)
					screen.Set(x, y, colors[i])
				}
			}
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three-Body Problem Simulation")
	ebiten.SetMaxTPS(120)

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
