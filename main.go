package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type PlayerBoard struct {
	xpos   float32
	ypos   float32
	speed  float32
	width  float32
	height float32
	color  rl.Color
}

func NewBoard(xpos, ypos, speed, width, height float32, color rl.Color) PlayerBoard {
	return PlayerBoard{
		xpos:   xpos,
		ypos:   ypos,
		speed:  speed,
		width:  width,
		height: height,
		color:  color,
	}
}

func (c PlayerBoard) DrawBoard() {
	rl.DrawRectangle(int32(c.xpos), int32(c.ypos), int32(c.width), int32(c.height), c.color)
}

func (c *PlayerBoard) PlayerMovement(xOffset float32) {
	c.xpos += xOffset

	if c.xpos < 0 {
		c.xpos = 0
	}
	if c.xpos+c.width > float32(rl.GetScreenWidth()) {
		c.xpos = float32(rl.GetScreenWidth()) - c.width
	}
}

func main() {
	rl.InitWindow(800, 400, "Weekly Project 6 - Breakout Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	//create the moving board
	player := NewBoard(400, 350, 10, 60, 5, rl.RayWhite)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkBlue)

		//draw the board player will be moving
		player.DrawBoard()

		if rl.IsKeyDown(rl.KeyA) {
			player.PlayerMovement(-player.speed)
		}
		if rl.IsKeyDown(rl.KeyD) {
			player.PlayerMovement(player.speed)
		}

		rl.EndDrawing()
	}
}
