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

// physics for the ball
type PhysicsBody struct {
	Pos    rl.Vector2
	Vel    rl.Vector2
	Mass   float32
	Radius float32
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

// function to draw the board using PlayerBoard struct
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

func NewBall(x, y float32) PhysicsBody {
	return PhysicsBody{
		Pos:    rl.NewVector2(x, y),
		Vel:    rl.NewVector2(0, 0),
		Mass:   1,
		Radius: 5,
	}
}

func (b PhysicsBody) DrawBall() {
	rl.DrawCircle(int32(b.Pos.X), int32(b.Pos.Y), b.Radius, rl.RayWhite)
}

func main() {
	rl.InitWindow(800, 400, "Weekly Project 6 - Breakout Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	//create the moving board
	player := NewBoard(400, 350, 5, 60, 5, rl.RayWhite)
	//initialize a ball
	ball := NewBall(player.xpos+player.width/2, player.ypos-10)

	ballLaunched := false

	//start game loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()

		rl.ClearBackground(rl.DarkBlue)

		//draw the board player will be moving
		player.DrawBoard()

		//movement keys for the paddle
		if rl.IsKeyDown(rl.KeyA) {
			player.PlayerMovement(-player.speed)
		}
		if rl.IsKeyDown(rl.KeyD) {
			player.PlayerMovement(player.speed)
		}

		//make sure the ball starts above the paddle board
		if !ballLaunched {
			ball.Pos.X = player.xpos + player.width/2
			ball.Pos.Y = player.ypos - ball.Radius
		}

		//use space to launch the static ball off the paddle
		if rl.IsKeyPressed(rl.KeySpace) && !ballLaunched {
			ball.Vel.Y = -3
			ball.Vel.X = player.speed * 0.5
			ballLaunched = true
		}

		//once ball is moving update its position
		ball.Pos.X += ball.Vel.X
		ball.Pos.Y += ball.Vel.Y

		//bounce off walls
		if ball.Pos.X-ball.Radius <= 0 || ball.Pos.X+ball.Radius >= float32(rl.GetScreenWidth()) {
			ball.Vel.X *= -1 //reverse the direction
		}
		if ball.Pos.Y-ball.Radius <= 0 {
			ball.Vel.Y *= -1 //reverse the direction if the ball hits the top
		}

		//check to see if ball falls below the screen
		if ball.Pos.Y > float32(rl.GetScreenHeight()) {
			ball = NewBall(player.xpos+player.width/2, player.ypos-10)
			ball.Vel = rl.NewVector2(0, 0) //reset velocity of ball
			ballLaunched = false
		}

		//bounce off the paddle
		if ball.Pos.Y+ball.Radius >= player.ypos &&
			ball.Pos.X > player.xpos &&
			ball.Pos.X < player.xpos+player.width {

			ball.Vel.Y = -1

			hitPosition := (ball.Pos.X - player.xpos) / player.width
			ball.Vel.X = (hitPosition - 0.5) * 6
		}

		//draw the ball
		ball.DrawBall()

		rl.EndDrawing()
	}
}
