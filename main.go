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

// struct for the blocks that will make a grid
type Block struct {
	Pos    rl.Vector2
	Width  float32
	Height float32
	Color  rl.Color
	Active bool //block has been hit or not
}

// physics for the ball
type PhysicsBody struct {
	Pos    rl.Vector2
	Vel    rl.Vector2
	Mass   float32
	Radius float32
}

// create a new paddle
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

// function to create a grid of blocks
func CreateBlocks(rows, cols int, blockWidth, blockHeight, spacing float32) []Block {
	var blocks []Block
	startX := float32(50) //top padding
	startY := float32(50) //bottom padding

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			x := startX + float32(col)*(blockWidth+spacing)
			y := startY + float32(row)*(blockHeight+spacing)
			block := Block{
				Pos:    rl.NewVector2(x, y),
				Width:  blockWidth,
				Height: blockHeight,
				Color:  rl.Orange,
				Active: true,
			}
			blocks = append(blocks, block)
		}
	}
	return blocks
}

// function to draw the blocks we made
func (b Block) DrawBlock() {
	if b.Active {
		rl.DrawRectangle(int32(b.Pos.X), int32(b.Pos.Y), int32(b.Width), int32(b.Height), b.Color)
	}
}

// function to check if ball collided with a square
func BallCollision(ball *PhysicsBody, block *Block) bool {
	if !block.Active {
		return false
	}

	//AABB collision check(axis-aligned bounding box)
	return ball.Pos.X+ball.Radius > block.Pos.X &&
		ball.Pos.X-ball.Radius < block.Pos.X+block.Width &&
		ball.Pos.Y+ball.Radius > block.Pos.Y &&
		ball.Pos.Y-ball.Radius < block.Pos.Y+block.Height
}

func main() {
	rl.InitWindow(800, 400, "Weekly Project 6 - Breakout Game")

	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	//create the moving board
	player := NewBoard(400, 350, 5, 60, 5, rl.RayWhite)
	//initialize a ball
	ball := NewBall(player.xpos+player.width/2, player.ypos-10)
	//draw the grid of blocks
	blocks := CreateBlocks(3, 11, 60, 50, 5)

	ballLaunched := false //if ball has been launched or not

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
			if rl.IsKeyDown(rl.KeyA) {
				ball.Vel.X = -2 // Move left if A is held
			} else if rl.IsKeyDown(rl.KeyD) {
				ball.Vel.X = 2 // Move right if D is held
			} else {
				ball.Vel.X = 0 // Go straight up if no key is held
			}
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
			blocks = CreateBlocks(3, 11, 60, 50, 5) // Reset blocks
		}

		//bounce off the paddle
		if ball.Pos.Y+ball.Radius >= player.ypos &&
			ball.Pos.X > player.xpos &&
			ball.Pos.X < player.xpos+player.width {

			ball.Vel.Y = -3

			hitPosition := (ball.Pos.X - player.xpos) / player.width
			ball.Vel.X = (hitPosition - 0.5) * 6
		}

		//draw the ball
		ball.DrawBall()

		//draw all the blocks
		for i := range blocks {
			blocks[i].DrawBlock()
		}

		//check through each block to see if it was collided with
		for i := range blocks {
			//blocks[i].DrawBlock()
			if !blocks[i].Active {
				continue
			}

			if BallCollision(&ball, &blocks[i]) {
				blocks[i].Active = false //remove block
				// //determine bounce direction
				block := blocks[i]

				//check which side the ball hit
				if ball.Pos.Y < block.Pos.Y { //hit from below
					ball.Vel.Y = -ball.Vel.Y
				} else if ball.Pos.Y > block.Pos.Y+block.Height { //hit from above
					ball.Vel.Y = -ball.Vel.Y
				} else if ball.Pos.X < block.Pos.X { //hit from left
					ball.Vel.X = -ball.Vel.X
				} else if ball.Pos.X > block.Pos.X+block.Width { //hit from right
					ball.Vel.X = -ball.Vel.X
				}
				//increase ball speed slightly
				ball.Vel.X *= 1.05
				ball.Vel.Y *= 1.05
			}
		}

		//reset game
		gameOver := true
		for _, block := range blocks {
			if block.Active {
				gameOver = false
				break
			}
		}

		if gameOver {
			ball = NewBall(player.xpos+player.width/2, player.ypos-10)
			ball.Vel = rl.NewVector2(0, 0)
			ballLaunched = false
			blocks = CreateBlocks(3, 11, 60, 50, 5)
		}

		rl.EndDrawing()
	}
}
