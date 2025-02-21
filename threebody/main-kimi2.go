package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480

	gravityConstant = 1.0         // 万有引力常数
	timeStep        = 1.0 / 180.0 // 时间步长
	tailFadeTime    = 30.0        // 尾迹消失时间 (秒)
)

type CelestialBody struct {
	mass         float64      // 质量
	position     [2]float64   // 位置 (x, y)
	velocity     [2]float64   // 速度 (vx, vy)
	acceleration [2]float64   // 加速度 (ax, ay)
	radius       float64      // 半径
	color        ebiten.Color // 颜色
	tail         []TailPoint  // 尾迹点
	tailTime     time.Time    // 尾迹最近更新时间
}

type TailPoint struct {
	position [2]float64 // 位置 (x, y)
	alpha    float64    // 透明度
}

type Simulation struct {
	bodies      []CelestialBody
	startTime   time.Time
	resetBuffer time.Duration // 重置缓冲时间
}

// 计算加速度
func (b *CelestialBody) calculateAcceleration(bodies []CelestialBody) {
	ax := 0.0
	ay := 0.0
	for _, other := range bodies {
		if other == b {
			continue
		}
		dx := other.position[0] - b.position[0]
		dy := other.position[1] - b.position[1]
		distance := math.Sqrt(dx*dx + dy*dy)
		if distance < b.radius+other.radius {
			// 碰撞检测
			// 可以在这里处理碰撞逻辑，例如重置模拟
			continue
		}
		force := gravityConstant * b.mass * other.mass / (math.Pow(distance, 3) + 1e-9) // 避免除以零
		ax += force * dx / other.mass
		ay += force * dy / other.mass
	}
	b.acceleration = [2]float64{ax, ay}
}

// RK4 方法更新位置和速度
func (b *CelestialBody) update(timeStep float64) {
	// RK4 方法
	var k1v, k2v, k3v, k4v [2]float64
	var k1r, k2r, k3r, k4r [2]float64

	// Calculate k1
	b.calculateAcceleration([]CelestialBody{})
	k1v = [2]float64{b.acceleration[0], b.acceleration[1]}
	k1r = [2]float64{b.velocity[0], b.velocity[1]}

	// Calculate k2
	tempVelocity := [2]float64{b.velocity[0] + k1v[0]*timeStep/2, b.velocity[1] + k1v[1]*timeStep/2}
	tempPosition := [2]float64{b.position[0] + k1r[0]*timeStep/2, b.position[1] + k1r[1]*timeStep/2}
	k2v = [2]float64{b.acceleration[0], b.acceleration[1]} // need to recalculate acceleration
	b.velocity = tempVelocity
	b.position = tempPosition
	b.calculateAcceleration([]CelestialBody{})
	k2v = [2]float64{b.acceleration[0], b.acceleration[1]}
	k2r = [2]float64{tempVelocity[0], tempVelocity[1]}

	// Calculate k3
	tempVelocity = [2]float64{b.velocity[0] + k2v[0]*timeStep/2, b.velocity[1] + k2v[1]*timeStep/2}
	tempPosition = [2]float64{b.position[0] + k2r[0]*timeStep/2, b.position[1] + k2r[1]*timeStep/2}
	k3v = [2]float64{b.acceleration[0], b.acceleration[1]}
	b.velocity = tempVelocity
	b.position = tempPosition
	b.calculateAcceleration([]CelestialBody{})
	k3v = [2]float64{b.acceleration[0], b.acceleration[1]}
	k3r = [2]float64{tempVelocity[0], tempVelocity[1]}

	// Calculate k4
	tempVelocity = [2]float64{b.velocity[0] + k3v[0]*timeStep, b.velocity[1] + k3v[1]*timeStep}
	tempPosition = [2]float64{b.position[0] + k3r[0]*timeStep, b.position[1] + k3r[1]*timeStep}
	k4v = [2]float64{b.acceleration[0], b.acceleration[1]}
	b.velocity = tempVelocity
	b.position = tempPosition
	b.calculateAcceleration([]CelestialBody{})
	k4v = [2]float64{b.acceleration[0], b.acceleration[1]}
	k4r = [2]float64{tempVelocity[0], tempVelocity[1]}

	// Update velocity and position
	b.velocity[0] += (k1v[0] + 2*k2v[0] + 2*k3v[0] + k4v[0]) * timeStep / 6
	b.velocity[1] += (k1v[1] + 2*k2v[1] + 2*k3v[1] + k4v[1]) * timeStep / 6
	b.position[0] += (k1r[0] + 2*k2r[0] + 2*k3r[0] + k4r[0]) * timeStep / 6
	b.position[1] += (k1r[1] + 2*k2r[1] + 2*k3r[1] + k4r[1]) * timeStep / 6

	// 检查边界碰撞
	border := 320.0 // 物理范围
	if b.position[0] < -border || b.position[0] > border || b.position[1] < -border || b.position[1] > border {
		// 位置超出边界，重置模拟
		sim := Simulation{bodies: []CelestialBody{
			{
				mass:     100.0,
				position: [2]float64{-100.0, 0.0},
				velocity: [2]float64{0.5, 1.0},
				radius:   5.0,
				color:    [4]float32{1.0, 0.0, 0.0, 1.0},
			},
			{
				mass:     100.0,
				position: [2]float64{100.0, 0.0},
				velocity: [2]float64{-0.5, -1.0},
				radius:   5.0,
				color:    [4]float32{0.0, 1.0, 0.0, 1.0},
			},
			{
				mass:     100.0,
				position: [2]float64{0.0, 100.0},
				velocity: [2]float64{1.0, 0.0},
				radius:   5.0,
				color:    [4]float32{0.0, 0.0, 1.0, 1.0},
			},
		}}
		fmt.Println("重置模拟")
		sim.Init()
		return &sim
	}

	// 添加尾迹点
	b.tail = append(b.tail, TailPoint{
		position: [2]float64{b.position[0], b.position[1]},
		alpha:    1.0,
	})
	b.tailTime = time.Now()

	// 更新尾迹透明度
	for i := range b.tail {
		elapsed := time.Since(b.tailTime).Seconds()
		fadeFactor := elapsed / tailFadeTime
		if fadeFactor > 1.0 {
			b.tail[i].alpha = 0
		} else {
			b.tail[i].alpha = 1.0 - fadeFactor
		}
	}

	// 清除超时尾迹
	b.tail = filterTail(b.tail)
	return
}

// 过滤尾迹，移除透明度为0的点
func filterTail(tail []TailPoint) []TailPoint {
	var newTail []TailPoint
	for _, p := range tail {
		if p.alpha > 0 {
			newTail = append(newTail, p)
		}
	}
	return newTail
}

func NewCelestialBody(mass float64, x, y, vx, vy float64, radius float64, color ebiten.Color) CelestialBody {
	return CelestialBody{
		mass:     mass,
		position: [2]float64{x, y},
		velocity: [2]float64{vx, vy},
		radius:   radius,
		color:    color,
	}
}

func (s *Simulation) Update() error {
	for _, body := range s.bodies {
		body.update(timeStep * time.Second)
	}
	return nil
}

func (s *Simulation) Draw(screen *ebiten.Image) {
	// 清屏
	screen.Fill(ebiten.ClearColor)

	// 绘制尾迹
	for _, body := range s.bodies {
		for _, p := range body.tail {
			if p.alpha <= 0 {
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(p.position[0]+float64(screenWidth)/2, p.position[1]+float64(screenHeight)/2)
			radius := body.radius / 2
			img := ebiten.NewImage(int(2*radius), int(2*radius))
			img.Fill(body.color)
			img.SetAllColor(ebiten.Color{A: uint8(p.alpha * 255)})
			screen.DrawImage(img, op)
		}
	}

	// 绘制球体
	for _, body := range s.bodies {
		op := &ebiten.DrawCircleOptions{}
		op.Color = body.color
		op.Radius = body.radius
		ebiten.DrawCircle(screen, body.position[0]+float64(screenWidth)/2, body.position[1]+float64(screenHeight)/2, op)
	}

	// 显示帧率
	fmt.Printf("FPS: %0.2f\n", ebiten.CurrentFPS())
}

func (s *Simulation) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (s *Simulation) Init() {
	// 初始化三个球体
	s.bodies = []CelestialBody{
		{
			mass:     100.0,
			position: [2]float64{-100.0, 0.0},
			velocity: [2]float64{0.5, 1.0},
			radius:   5.0,
			color:    [4]float32{1.0, 0.0, 0.0, 1.0},
		},
		{
			mass:     100.0,
			position: [2]float64{100.0, 0.0},
			velocity: [2]float64{-0.5, -1.0},
			radius:   5.0,
			color:    [4]float32{0.0, 1.0, 0.0, 1.0},
		},
		{
			mass:     100.0,
			position: [2]float64{0.0, 100.0},
			velocity: [2]float64{1.0, 0.0},
			radius:   5.0,
			color:    [4]float32{0.0, 0.0, 1.0, 1.0},
		},
	}
	s.startTime = time.Now()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ebiten.SetWindowTitle("三体问题模拟")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetFullscreen(false)
	ebiten.SetMaxTPS(120)

	sim := Simulation{}
	sim.Init()

	if err := ebiten.RunGame(&sim); err != nil {
		fmt.Println(err)
	}
}
