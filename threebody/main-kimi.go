package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 640
	screenHeight = 480
	maxSpeed     = 5
)

type Ball struct {
	x, y, vx, vy, mass float64
	color              color.Color
}

func NewBall() *Ball {
	return &Ball{
		x:     float64(rand.Intn(screenWidth)),
		y:     float64(rand.Intn(screenHeight)),
		vx:    (rand.Float64()*2 - 1) * maxSpeed,
		vy:    (rand.Float64()*2 - 1) * maxSpeed,
		mass:  rand.Float64()*100 + 0.1,
		color: color.RGBA{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255},
	}
}

func (b *Ball) Update(balls []*Ball) {
	ax, ay := 0.0, 0.0
	for _, other := range balls {
		if other == b {
			continue
		}
		dx := other.x - b.x
		dy := other.y - b.y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 5 { // 防止除零错误，同时限制最小距离
			dist = 5
		}
		g := (100 * other.mass) / (dist * dist) // 引力公式（简化版本）
		ax += g * dx / dist
		ay += g * dy / dist
	}
	b.vx += ax
	b.vy += ay
	b.x += b.vx
	b.y += b.vy
	// 边界检测
	if b.x < 0 || b.x > screenWidth || b.y < 0 || b.y > screenHeight {
		InitBalls()
	}
}

func DrawBall(screen *ebiten.Image, b *Ball) {
	screen.Fill(color.Black) // 背景颜色
	radius := 20
	for _, ball := range allBalls {
		ebiten.DrawCircle(screen, int(ball.x), int(ball.y), radius, ball.color)
	}
}

var allBalls []*Ball

func InitBalls() {
	allBalls = make([]*Ball, 0, 3)
	for i := 0; i < 3; i++ {
		allBalls = append(allBalls, NewBall())
	}
}

func Update(screen *ebiten.Image) error {
	for _, ball := range allBalls {
		ball.Update(allBalls)
	}
	DrawBall(screen, nil)
	return nil
}

func main() {
	rand.Seed(0)
	InitBalls()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three Body Problem Simulation")
	ebiten.Run(Update, screenWidth, screenHeight, 60, "Three Body Problem Simulation")
}
