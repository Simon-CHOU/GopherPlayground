package main

import (
    "log"
    "time"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/text"
    "golang.org/x/image/font/sansmath"
)

const (
    screenWidth  = 800
    screenHeight = 600
    g           = 0.1 // 重力常数，此处简化处理
)

type Ball struct {
    x, y      float64
    vx, vy    float64
    r         float64
    mass      float64
}

func (b *Ball) Draw(screen *ebiten.Image) {
    // 绘制球体
    ebiten.DrawCircle(screen, b.x, b.y, b.r, ebiten.ColorMatcher{R: 255, G: 0, B: 0, A: 255})
}

func main() {
    // 初始化三个球体，每个有不同的初始位置、速度和半径
    balls := []Ball{
        {x: 100, y: 300, vx: 10, vy: 0, r: 20, mass: 1},
        {x: 400, y: 300, vx: -5, vy: 0, r: 20, mass: 2},
        {x: 700, y: 300, vx: 0, vy: 10, r: 20, mass: 3},
    }

    // 初始化Ebiten窗口
    screen := ebiten.NewImage(screenWidth, screenHeight)
    ebiten.SetWindowSize(screenWidth, screenHeight)

    for {
        // 清除屏幕（绘制背景）
        screen.Clear()
        drawBackground(screen)

        // 更新每个球体的位置和速度
        for i := range balls {
            balls[i].Update(&balls)
        }

        // 绘制所有球体
        for _, b := range balls {
            b.Draw(screen)
        }

        // 显示屏幕内容并处理输入事件
        if err := ebiten.Update(func(g *ebiten.Game) error {
            ebitenutil.DebugDraw(g Screen, screen, 0, 0)
            return nil
        }); err != nil {
            log.Fatal(err)
        }

        // 处理窗口关闭事件
        select {
        case <-time.After(10 * time.Millisecond):
        case e := <-ebiten.Events:
            if e.Key == ebiten.KeyQuit {
                return
            }
        }
    }
}

func drawBackground(screen *ebiten.Image) {
    // 绘制背景网格（可选）
    text.Draw(screen, " ", sansmath.Tiny, 0, 0, ebiten.ColorMatcher{R: 255, G: 255, B: 255, A: 255})
}

func (b *Ball) Update(balls []Ball) {
    // 处理边界条件，当球体超出屏幕时重新开始模拟
    if b.x < -b.r || b.x > screenWidth+b.r || b.y < -b.r || b.y > screenHeight+b.r {
        ResetBall(b)
        return
    }

    // 计算与其他球体的引力并更新速度和位置
    for j := range balls {
        if i == j {
            continue
        }
        others := []int{j}
        for k := range balls {
            if k != i && k != j {
                others = append(others, k)
            }
        }
        // 简化计算，只考虑与其他两个球体的引力（实际应考虑所有其他球体）
        jBall := &balls[j]
        dx := b.x - jBall.x
        dy := b.y - jBall.y
        distanceSquared := dx*dx + dy*dy
        if distanceSquared == 0 {
            continue // 避免同一位置的情况
        }
        force := g * (b.mass * jBall.mass) / distanceSquared
        angle := math.Atan2(dy, dx)
        b.vx -= force * dx / distance * b.mass
        b.vy -= force * dy / distance * b.mass
    }

    // 更新位置
    b.x += b.vx
    b.y += b.vy

    // 简化速度衰减（可选）
    b.vx *= 0.99
    b.vy *= 0.99
}

func ResetBall(b *Ball) {
    // 随机重置球体的位置和速度
    r := 50 // 半径
    angle := rand.Float64() * math.Pi * 2
    x := screenWidth/2 + r*math.Cos(angle)
    y := screenHeight/2 + r*math.Sin(angle)
    speed := 3
    b.x = x
    b.y = y
    b.vx = speed * math.Cos(angle)
    b.vy = speed * math.Sin(angle)
}