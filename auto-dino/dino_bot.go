package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/go-vgo/robotgo"
)

const (
	// Dino game constants
	jumpKey         = "space"
	retryKey        = "space"
	gameOverColor   = "535353"
	initialDelay    = 3 * time.Second
	retryDelay      = 2 * time.Second
	screenWidth     = 1920
	screenHeight    = 1080
	gameWindowWidth = 600
	gameWindowHeight = 150
)

func main() {
	fmt.Println("Starting Dino Bot...")
	time.Sleep(initialDelay)

	for {
		// Capture the game window
		gameWindow := robotgo.CaptureScreen((screenWidth-gameWindowWidth)/2, (screenHeight-gameWindowHeight)/2, gameWindowWidth, gameWindowHeight)
		if gameWindow == nil {
			log.Fatal("Failed to capture the game window")
		}
		
		img := robotgo.ToImage(gameWindow)
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
	for x := img.Bounds().Min.X + 100; x < img.Bounds().Min.X + 150; x++ {
		for y := img.Bounds().Min.Y + 50; y < img.Bounds().Min.Y + 100; y++ {
			if img.At(x, y) == obstacleColor {
				return true
			}
		}
	}
	return false
}