package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	// Dino game constants
	jumpKey          = "space"
	retryKey         = "space"
	initialDelay     = 3 * time.Second
	retryDelay       = 2 * time.Second
	gameWindowWidth  = 600
	gameWindowHeight = 150
)

func main() {

	// Detect screen size
	screenWidth, screenHeight := robotgo.GetScreenSize()

	// Determine the primary monitor (assuming the primary monitor has larger width)
	primaryMonitorWidth := screenWidth / 2
	primaryMonitorHeight := screenHeight

	// Positions to place the game window in the center of the primary monitor
	gameWindowX := (primaryMonitorWidth - gameWindowWidth) / 2
	gameWindowY := (primaryMonitorHeight - gameWindowHeight) / 2

	fmt.Println("Starting Dino Bot...")
	drawGuideFrame(gameWindowX, gameWindowY, gameWindowWidth, gameWindowHeight)
	time.Sleep(initialDelay)

	for {
		// Capture the game window
		gameWindow := robotgo.CaptureScreen(gameWindowX, gameWindowY, gameWindowWidth, gameWindowHeight)
		if gameWindow == nil {
			log.Fatal("Failed to capture the game window")
		}

		img := robotgo.ToImage(gameWindow)
		// 预览截图
		// 保存图片到 ./img/ 目录
		if err := saveImage(&img, "./img"); err != nil {
			log.Printf("Failed to save image: %v", err)
		}

		robotgo.FreeBitmap(gameWindow)

		// Check for obstacles and game over
		if isGameOver(img) {
			log.Println("Game over detected. Restarting...")
			robotgo.KeyTap(retryKey)
			time.Sleep(retryDelay)
			continue
		}

		if isObstacleDetected(img) {
			log.Println("Obstacle detected. Jumping...")
			robotgo.KeyTap(jumpKey)
		}

		time.Sleep(50 * time.Millisecond)
	}
}

func drawGuideFrame(x, y, width, height int) {
	fmt.Printf("\nPlease place the Dino game window within the following coordinates:\n")
	fmt.Printf("Top-left corner: (%d, %d)\n", x, y)
	fmt.Printf("Bottom-right corner: (%d, %d)\n", x+width, y+height)
	fmt.Println("Press Enter to start the bot...")
	var input string
	fmt.Scanln(&input)
}

func isGameOver(img image.Image) bool {
	// Check for game over color in the middle of the screen
	gameOverColor := color.RGBA{R: 83, G: 83, B: 83, A: 255}
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
			if img.At(x, y) == gameOverColor {
				return true
			}
		}
	}
	return false
}

func isObstacleDetected(img image.Image) bool {
	// Check for obstacles in front of the dinosaur
	obstacleColor := color.RGBA{R: 83, G: 83, B: 83, A: 255}
	for x := img.Bounds().Min.X + 100; x < img.Bounds().Min.X+150; x++ {
		for y := img.Bounds().Min.Y + 50; y < img.Bounds().Min.Y+100; y++ {
			if img.At(x, y) == obstacleColor {
				return true
			}
		}
	}
	return false
}

// 保存图片

func saveImage(img *image.Image, dir string) error {

	// 获取当前时间
	now := time.Now()

	// 格式化日期时间
	dateStr := now.Format("200601021504")

	// 获取三位数序号
	var i int
	for {
		filename := fmt.Sprintf("%s/%s%03d.png", dir, dateStr, i)
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// 文件不存在，可以使用该序号
			err := saveImageToFile(img, filename)
			if err != nil {
				return fmt.Errorf("保存图片失败: %w", err)
			}
			fmt.Println("图片已保存:", filename)
			return nil
		}
		i++
	}
}

func saveImageToFile(img *image.Image, filename string) error {
	// 创建目录，如果不存在
	err := os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return err
	}

	// 创建文件
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// 保存图像为 PNG 格式
	err = png.Encode(f, *img)
	if err != nil {
		return err
	}

	return nil
}
