package main

import (
	"fmt"
	"math/rand"

	"github.com/mattn/go-tty"
)

type point struct {
	x, y int
}

type snake struct {
	body []point
}

type game struct {
	apple         point
	height, width int
	over          bool
	score         int
	snake         snake
}

func (p *point) same(point point) bool {
	return p.x == point.x && p.y == point.y
}

// Grow the snake by 1 point
// by appending a point to the end
// of snake.body.
func (s *snake) grow() {
	s.body = append(s.body, point{})
}

// Move the snake by (x, y)
func (s *snake) move(x, y int) {
	// Loop backwards through the snake's body
	// and move each point to the position of
	// the point before it.
	//
	// Also, move the snake's body before
	// moving the head so that the head
	// doesn't overwrite the position of
	// the point before it.
	for i := len(s.body) - 1; i > 0; i-- {
		s.body[i] = s.body[i-1]
	}

	// Move the snake's head.
	s.body[0].x += x
	s.body[0].y += y
}

func (g *game) input() {
	tty, _ := tty.Open()
	defer tty.Close()

	ch, err := tty.ReadRune()
	if err != nil {
		panic(err)
	}
	switch ch {
	case 'w', 'k':
		g.snake.move(0, -1)
	case 'a', 'h':
		g.snake.move(-1, 0)
	case 's', 'j':
		g.snake.move(0, 1)
	case 'd', 'l':
		g.snake.move(1, 0)
	case 'q':
		g.over = true
	}
}

func (g *game) update() {
	// Check if snake overlaps with apple.
	if g.apple.same(g.snake.body[0]) {
		// Increase score by 1.
		g.score++

		// Grow snake.
		g.snake.grow()

		// Create a new point for the apple
		// at a random position.
		g.apple = point{rand.Intn(g.width), rand.Intn(g.height)}
	}

	// Check if snake overlaps with itself.
	for _, p := range g.snake.body[1:] {
		if g.snake.body[0].same(p) {
			g.over = true
		}
	}

	// Wrap snake around screen.
	head := &g.snake.body[0]
	head.x = (head.x + g.width) % g.width
	head.y = (head.y + g.height) % g.height
}

func (g *game) draw() {
	// Clear the console.
	fmt.Print("\033[H\033[2J")

	// Print the score.
	fmt.Printf("Score: %d\n", g.score)

	h := g.height
	w := g.width
	buf := make([]byte, 0, h*(w))

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			switch {
			case g.snake.body[0].same(point{x, y}):
				buf = append(buf, '@')
			case g.apple.same(point{x, y}):
				buf = append(buf, '$')
			default:
				buf = append(buf, '.')
			}
		}
		buf = append(buf, '\n')
	}

	// For every point on the snake's body,
	// after the first point, print an 'o'.
	//
	// i is calculated by multiplying the
	// y coordinate by the width of the
	// screen and adding the x coordinate.
	for _, p := range g.snake.body[1:] {
		i := p.y*(w+1) + p.x
		buf[i] = 'o'
	}

	fmt.Print(string(buf))

	// Print game controls.
	fmt.Println("wasd or hjkl to move, q to quit")
}

func main() {
	g := game{
		height: 10,
		width:  20,
		// Initialize snake at (0, 0).
		snake: snake{
			body: []point{
				{rand.Intn(20), rand.Intn(10)},
			},
		},
		// Initialize apple at a random position.
		apple: point{rand.Intn(20), rand.Intn(10)},
	}

	for !g.over {
		g.draw()
		g.input()
		g.update()
	}

	// Show the cursor.
	fmt.Print("\033[?25h")

	// Print game over message and final score.
	fmt.Printf("Game Over! Score: %d\n", g.score)
}
