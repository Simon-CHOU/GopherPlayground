package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 800
	screenHeight = 600
	G            = 1200.0 // 引力常数
	dt           = 0.008  // 时间步长
	minDistance  = 50.0   // 最小有效距离
)

type Body struct {
	posX, posY float64
	velX, velY float64
	mass       float64
	radius     float64
	color      color.Color
}

type Game struct {
	bodies []*Body
}

func (g *Game) Update() error {
	accelerations := make([]struct{ ax, ay float64 }, len(g.bodies))

	// 计算加速度
	for i := range g.bodies {
		var ax, ay float64
		for j := range g.bodies {
			if i == j {
				continue
			}

			dx := g.bodies[j].posX - g.bodies[i].posX
			dy := g.bodies[j].posY - g.bodies[i].posY
			rSq := dx*dx + dy*dy

			// 添加距离限制以避免数值不稳定
			if rSq < minDistance*minDistance {
				rSq = minDistance * minDistance
			}

			r := math.Sqrt(rSq)
			accel := G * g.bodies[j].mass / rSq
			ax += accel * dx / r
			ay += accel * dy / r
		}
		accelerations[i] = struct{ ax, ay float64 }{ax, ay}
	}

	// 更新速度和位置
	for i, body := range g.bodies {
		acc := accelerations[i]
		body.velX += acc.ax * dt
		body.velY += acc.ay * dt
		body.posX += body.velX * dt
		body.posY += body.velY * dt
	}

	// 边界检查
	for _, body := range g.bodies {
		if body.posX < 0 || body.posX > screenWidth ||
			body.posY < 0 || body.posY > screenHeight {
			g.reset()
			break
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 25, G: 25, B: 25, A: 255}) // 深色背景
	for _, body := range g.bodies {
		vector.DrawFilledCircle(screen,
			float32(body.posX),
			float32(body.posY),
			float32(body.radius),
			body.color,
			true)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) reset() {
	g.bodies = []*Body{
		{
			posX:   300,                        // 初始X位置
			posY:   300,                        // 初始Y位置
			velX:   0,                          // 初始X速度
			velY:   80,                         // 初始Y速度
			mass:   1000,                       // 质量
			radius: 12,                         // 显示半径
			color:  color.RGBA{R: 255, A: 255}, // 红色
		},
		{
			posX:   500,
			posY:   300,
			velX:   0,
			velY:   -80,
			mass:   1000,
			radius: 12,
			color:  color.RGBA{G: 255, A: 255}, // 绿色
		},
		{
			posX:   400,
			posY:   200,
			velX:   80,
			velY:   0,
			mass:   1000,
			radius: 12,
			color:  color.RGBA{B: 255, A: 255}, // 蓝色
		},
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("三体运动模拟")
	game := &Game{}
	game.reset()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
