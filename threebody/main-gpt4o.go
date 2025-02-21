package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

// 窗口大小
const (
	screenWidth  = 800
	screenHeight = 600
	// dt 为时间步长，可以调节精度与运动效果
	dt = 0.1
	// G 为引力常数，值可调节
	G = 0.1
)

// Ball 表示一个运动的球体（模拟行星）
type Ball struct {
	x, y   float64    // 位置
	vx, vy float64    // 速度
	mass   float64    // 质量
	radius float64    // 半径
	col    color.RGBA // 颜色
}

// Game 保存球体列表和预先绘制好的圆形图像
type Game struct {
	balls       []*Ball
	circleImage *ebiten.Image
}

// NewGame 初始化一个 Game，并生成球体及圆形图像
func NewGame() *Game {
	g := &Game{
		balls: make([]*Ball, 3),
	}
	g.reset()
	// 创建一个直径为球直径的图像（此处球半径为 10）
	g.circleImage = createCircleImage(int(10))
	return g
}

// reset 用固定的初始条件重新初始化球体
func (g *Game) reset() {
	// 为了增加随机性，也可以使用 rand 生成初始值
	rand.Seed(time.Now().UnixNano())
	g.balls[0] = &Ball{
		x:      300,
		y:      300,
		vx:     0,
		vy:     20,
		mass:   10,
		radius: 10,
		col:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
	}
	g.balls[1] = &Ball{
		x:      500,
		y:      300,
		vx:     0,
		vy:     -20,
		mass:   10,
		radius: 10,
		col:    color.RGBA{R: 0, G: 255, B: 0, A: 255},
	}
	g.balls[2] = &Ball{
		x:      400,
		y:      200,
		vx:     20,
		vy:     0,
		mass:   10,
		radius: 10,
		col:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}
}

// createCircleImage 生成一个以白色绘制的实心圆图像，半径为 r，尺寸为 r*2
// 绘制后通过 ColorM.Scale 可以轻松变换为其他颜色
func createCircleImage(r int) *ebiten.Image {
	diameter := r * 2
	img := ebiten.NewImage(diameter, diameter)
	// 将整个图像设置为透明背景
	img.Fill(color.Transparent)
	// 对每个像素点判断是否在圆内
	for y := 0; y < diameter; y++ {
		for x := 0; x < diameter; x++ {
			dx := float64(x - r)
			dy := float64(y - r)
			if dx*dx+dy*dy <= float64(r*r) {
				img.Set(x, y, color.White)
			}
		}
	}
	return img
}

// Update 每帧更新物理状态：计算引力、更新速度与位置
func (g *Game) Update() error {
	// 对每个球体计算因其他球体产生的引力加速度
	for i, ball := range g.balls {
		ax, ay := 0.0, 0.0
		for j, other := range g.balls {
			if i == j {
				continue
			}
			dx := other.x - ball.x
			dy := other.y - ball.y
			distSq := dx*dx + dy*dy
			// 防止距离过小导致除零错误
			if distSq < 1 {
				distSq = 1
			}
			// 计算牛顿万有引力（F = G*m1*m2/d^2），进而得到加速度（a = F/m1 = G*m2/d^2）
			a := G * other.mass / distSq
			dist := math.Sqrt(distSq)
			ax += a * dx / dist
			ay += a * dy / dist
		}
		// 更新速度
		ball.vx += ax * dt
		ball.vy += ay * dt
	}
	// 更新位置
	for _, ball := range g.balls {
		ball.x += ball.vx * dt
		ball.y += ball.vy * dt
	}
	// 检查是否有球体超出屏幕边界，如有则重置模拟
	for _, ball := range g.balls {
		if ball.x-ball.radius < 0 || ball.x+ball.radius > screenWidth ||
			ball.y-ball.radius < 0 || ball.y+ball.radius > screenHeight {
			g.reset()
			break
		}
	}
	return nil
}

// Draw 将球体绘制到屏幕上
func (g *Game) Draw(screen *ebiten.Image) {
	// 清空屏幕背景为黑色
	screen.Fill(color.Black)
	// 对每个球体绘制预先创建好的圆形图像，并根据球体颜色着色
	for _, ball := range g.balls {
		op := &ebiten.DrawImageOptions{}
		// 将图像平移到球心位置（注意图像中心在 (radius, radius)）
		op.GeoM.Translate(ball.x-ball.radius, ball.y-ball.radius)
		// 由于 circleImage 是白色图像，通过 ColorM.Scale 可以改变颜色
		op.ColorM.Scale(
			float64(ball.col.R)/255,
			float64(ball.col.G)/255,
			float64(ball.col.B)/255,
			1,
		)
		screen.DrawImage(g.circleImage, op)
	}
}

// Layout 设置逻辑屏幕尺寸
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("三体问题模拟")
	game := NewGame()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
