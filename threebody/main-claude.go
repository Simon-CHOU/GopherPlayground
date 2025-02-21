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
	G            = 1.0  // 增大引力常数使运动更快
	radius       = 10.0 // 球体半径
	tailLength   = 200  // 增加尾迹长度
	tailLife     = 30.0 // 尾迹生命周期（秒）
)

type Body struct {
	x, y     float64   // 位置
	vx, vy   float64   // 速度
	mass     float64   // 质量
	tail     []Point   // 尾迹
	tailTime []float64 // 尾迹产生的时间
}

type Point struct {
	x, y float64
}

type Game struct {
	bodies [3]*Body
	start  time.Time
}

func NewGame() *Game {
	// 初始化三个天体，增加初始速度
	bodies := [3]*Body{
		{
			x: 200, y: 300,
			vx: 2.0, vy: 1.0, // 增加初始速度
			mass:     1,
			tail:     make([]Point, 0, tailLength),
			tailTime: make([]float64, 0, tailLength),
		},
		{
			x: 400, y: 300,
			vx: -1.5, vy: 1.5, // 增加初始速度
			mass:     1,
			tail:     make([]Point, 0, tailLength),
			tailTime: make([]float64, 0, tailLength),
		},
		{
			x: 600, y: 300,
			vx: -1.0, vy: -1.5, // 增加初始速度
			mass:     1,
			tail:     make([]Point, 0, tailLength),
			tailTime: make([]float64, 0, tailLength),
		},
	}

	return &Game{
		bodies: bodies,
		start:  time.Now(),
	}
}

func (g *Game) Update() error {
	dt := 1.0 / 120.0 // 时间步长调整为120fps

	// 更新位置和速度
	for i := range g.bodies {
		// 计算引力
		fx, fy := 0.0, 0.0
		for j := range g.bodies {
			if i != j {
				dx := g.bodies[j].x - g.bodies[i].x
				dy := g.bodies[j].y - g.bodies[i].y
				r := math.Sqrt(dx*dx + dy*dy)

				// 检测碰撞
				if r < radius*2 {
					g.reset()
					return nil
				}

				f := G * g.bodies[i].mass * g.bodies[j].mass / (r * r)
				fx += f * dx / r
				fy += f * dy / r
			}
		}

		// 更新速度
		g.bodies[i].vx += fx * dt / g.bodies[i].mass
		g.bodies[i].vy += fy * dt / g.bodies[i].mass

		// 更新位置
		g.bodies[i].x += g.bodies[i].vx * dt
		g.bodies[i].y += g.bodies[i].vy * dt

		// 检查是否超出边界
		if g.bodies[i].x < 0 || g.bodies[i].x > screenWidth ||
			g.bodies[i].y < 0 || g.bodies[i].y > screenHeight {
			g.reset()
			return nil
		}

		// 更新尾迹
		now := time.Now()
		g.bodies[i].tail = append(g.bodies[i].tail, Point{g.bodies[i].x, g.bodies[i].y})
		g.bodies[i].tailTime = append(g.bodies[i].tailTime, now.Sub(g.start).Seconds())

		// 移除过期的尾迹点
		currentTime := now.Sub(g.start).Seconds()
		for len(g.bodies[i].tail) > 0 && currentTime-g.bodies[i].tailTime[0] > tailLife {
			g.bodies[i].tail = g.bodies[i].tail[1:]
			g.bodies[i].tailTime = g.bodies[i].tailTime[1:]
		}

		// 限制尾迹长度
		if len(g.bodies[i].tail) > tailLength {
			g.bodies[i].tail = g.bodies[i].tail[len(g.bodies[i].tail)-tailLength:]
			g.bodies[i].tailTime = g.bodies[i].tailTime[len(g.bodies[i].tailTime)-tailLength:]
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 绘制尾迹
	for i, body := range g.bodies {
		currentTime := time.Now().Sub(g.start).Seconds()

		// 为每个球体设置不同的基础颜色
		baseColors := []color.RGBA{
			{R: 255, G: 50, B: 50, A: 255}, // 红色
			{G: 255, B: 50, A: 255},        // 绿色
			{R: 50, G: 50, B: 255, A: 255}, // 蓝色
		}

		// 绘制尾迹
		for j := range body.tail {
			age := currentTime - body.tailTime[j]
			if age > tailLife {
				continue
			}

			// 计算透明度，使用平滑的渐变效果
			alpha := 1.0 - (age / tailLife)
			alpha = math.Pow(alpha, 0.5) // 使用平方根使渐变更加平滑

			col := color.RGBA{
				R: baseColors[i].R,
				G: baseColors[i].G,
				B: baseColors[i].B,
				A: uint8(255 * alpha),
			}

			// 尾迹粒子大小随时间渐变
			particleSize := 2.0 + 3.0*(1.0-alpha)
			DrawCircle(screen, body.tail[j].x, body.tail[j].y, particleSize, col)
		}

		// 绘制球体
		DrawCircle(screen, body.x, body.y, radius, baseColors[i])
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) reset() {
	*g = *NewGame()
}

// DrawCircle 绘制实心圆
func DrawCircle(screen *ebiten.Image, x, y, radius float64, clr color.Color) {
	r := int(radius)
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				screen.Set(int(x)+dx, int(y)+dy, clr)
			}
		}
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Three Body Simulation")
	ebiten.SetTPS(120)    // 设置刷新率为120fps
	ebiten.SetMaxTPS(120) // 设置最大刷新率为120fps

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
