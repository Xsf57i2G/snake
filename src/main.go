/*
 * The game of Snake.
 */

package main

import (
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
	apple point
	h, w  int
	over  bool
	score int
	snake snake
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

	input, _ := tty.ReadRune()

	switch input {
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

	tty.Close()
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
		h := rand.Intn(g.h)
		w := rand.Intn(g.w)
		g.apple = point{w, h}
	}

	// Check if snake overlaps with itself.
	for _, p := range g.snake.body[1:] {
		if g.snake.body[0].same(p) {
			g.over = true
		}
	}

	// Wrap snake around screen.
	head := &g.snake.body[0]
	head.x = (head.x + g.w) % g.w
	head.y = (head.y + g.h) % g.h
}

func (g *game) draw() {
	// Clear the console.
	print("\033[H\033[2J")

	// Print the score.
	println("Score:", g.score)

	var buf = make([]byte, 0, g.h*(g.w+1))

	for y := 0; y < g.h; y++ {
		for x := 0; x < g.w; x++ {
			switch {
			case g.snake.body[0].same(point{x, y}):
				buf = append(buf, '@')
			case g.apple.same(point{x, y}):
				buf = append(buf, '*')
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
		y := p.y
		x := p.x
		w := g.w
		i := y*(w+1) + x
		buf[i] = 'o'
	}

	println(string(buf))

	// Print game controls.
	println("wasd or hjkl to move, q to quit")
}

func main() {
	game := new(game)

	game.h = 10
	game.w = 20

	game.snake = snake{[]point{{game.w / 2, game.h / 2}}}
	game.apple = point{rand.Int() % game.w, rand.Int() % game.h}

	for !game.over {
		game.draw()
		game.input()
		game.update()
	}

	// Show the cursor.
	print("\033[?25h")

	// Print game over message and final score.
	print("Game Over! Score:", game.score)
}
