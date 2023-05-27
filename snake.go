package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const (
	width  = 40
	height = 20
	UP     = 'k'
	DOWN   = 'j'
	LEFT   = 'h'
	RIGHT  = 'l'
)

func contains(points []point, target point) bool {
	for _, point := range points {
		if point.x == target.x && point.y == target.y {
			return true
		}
	}
	return false
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

type point struct {
	x, y int
}

type gameState struct {
	coordination []point
	fruit        point
	direction    string
	gameOver     bool
	score        int
}

func (g *gameState) drawGame() {
	fmt.Println("Score: ", g.score)
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			if i == 0 || i == width-1 || j == 0 || j == height-1 {
				fmt.Print("#")
			} else if g.coordination[0].x == i && g.coordination[0].y == j {
				if g.gameOver {
					fmt.Print("X")
				} else {
					fmt.Print("0")
				}
			} else if contains(g.coordination[1:], point{i, j}) { // we skip the head
				fmt.Print("o")
			} else if g.fruit.x == i && g.fruit.y == j {
				fmt.Print("F")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

}

func (state *gameState) handleInput() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		key, _, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		if key == UP && state.direction != "down" {
			state.direction = "up"
		} else if key == DOWN && state.direction != "up" {
			state.direction = "down"
		} else if key == LEFT && state.direction != "right" {
			state.direction = "left"
		} else if key == RIGHT && state.direction != "left" {
			state.direction = "right"
		} else if key == 'q' || key == 'Q' {
			state.gameOver = true
			return
		}
	}
}

func (state *gameState) update() {
	newHead := state.coordination[0]
	if state.direction == "up" {
		newHead.y -= 1
	} else if state.direction == "down" {
		newHead.y += 1
	} else if state.direction == "left" {
		newHead.x -= 2
	} else if state.direction == "right" {
		newHead.x += 2
	}

	if contains(state.coordination, newHead) ||
		newHead.x <= 0 || newHead.x >= width-1 ||
		newHead.y <= 0 || newHead.y >= height-1 {
		state.gameOver = true
		return
	}
	state.coordination = append([]point{newHead}, state.coordination...)

	if state.fruitCollision(newHead) {
		state.score += 1
		state.generateFruit()
	} else { // if the snake doesn't eat the fruit we cut off the tail by one
		// to simulate movement
		state.coordination = state.coordination[:len(state.coordination)-1]
	}
}

func (state *gameState) checkState() {
	if state.gameOver {
		fmt.Println("Game Over")
		os.Exit(0)
	}
}

func (state *gameState) fruitCollision(newHead point) bool {
	dir := state.direction
	head := state.coordination[0]
	fruit := state.fruit
	// we do this to avoid the snake going through the fruit (snake moves 2 units horizontally)
	return dir == "left" && (head.x-1 == fruit.x && head.y == fruit.y) ||
		dir == "right" && (head.x+1 == fruit.x && head.y == fruit.y) ||
		newHead.x == state.fruit.x && newHead.y == state.fruit.y

}

func (state *gameState) generateFruit() {
	rand.Seed(time.Now().Unix())

	// starting width for head is width/2, so we need to make sure
	// that the fruit is eatable since we move by 2 units horizontally
	w := width / 2 % 2

	// we increment by one to avoid the border
	randomWidth := rand.Intn(width-3) + 1
	randomHeight := rand.Intn(height-3) + 1
	for randomWidth%2 != w {
		randomWidth = rand.Intn(width - 3)
	}
	state.fruit = point{randomWidth, randomHeight}
}

func main() {
	state := gameState{
		coordination: []point{point{width / 2, height / 2}},
		direction:    "right",
	}
	state.generateFruit()
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	go state.handleInput()

	for range ticker.C {
		select {
		default:
			clearScreen()
			state.update()
			state.drawGame()
			state.checkState()
		}
	}
}
