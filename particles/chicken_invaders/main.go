package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	w = 70
	h = 40
)

var jett Jett
var screen [][]rune

type Jett struct {
	x, y float32
	char rune
}

}

}


func (j *Jett) goRight() {
	j.x += 1
}

func (j *Jett) goLeft() {
	j.x -= 1
}

func (j *Jett) goUp() {
	j.y -= 1
}

func (j *Jett) goDown() {
	j.y += 1
}

func CreateJett() Jett {
	return Jett{
		x:    w / 2,
		y:    h - 2,
		char: '█',
	}
}

}

func ClearScreen() {
	screen = make([][]rune, h)
	for i := range screen {
		screen[i] = make([]rune, w)
		for j := range screen[i] {
			if i == 0 || i == h-1 {
				screen[i][j] = '─'
			} else if j == 0 || j == w-1 {
				screen[i][j] = '│'
			} else {
				screen[i][j] = ' '
			}
		}
	}

	screen[0][0] = '┌'
	screen[0][w-1] = '┐'
	screen[h-1][0] = '└'
	screen[h-1][w-1] = '┘'
}

func UpdateScreen() {
	if int(jett.x) >= 1 && int(jett.x) <= w-1 && int(jett.y) >= 1 && int(jett.y) <= h-1 {
		screen[int(jett.y)][int(jett.x)] = jett.char
	}

		}

	}
}

func handleInput(stop chan struct{}) {
	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			return
		}
		if buf[0] == 27 && buf[1] == '[' && n >= 3 {
			switch buf[2] {
			case 'A':
				jett.goUp()
			case 'B':
				jett.goDown()
			case 'C':
				jett.goRight()
			case 'D':
				jett.goLeft()
			}
		} else if buf[0] == 'q' || buf[0] == 27 {
			close(stop)
			return
		}
	}
}

func Render(stop chan struct{}) {
	for {
		select {
		case <-stop:
			return
		default:
		}

		ClearScreen()
		UpdateScreen()

		fmt.Print("\033[H\033[2J")
		for _, row := range screen {
			fmt.Print(string(row) + "\r\n")
		}

		time.Sleep(120 * time.Millisecond)
	}
}

func main() {

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	jett = CreateJett()

	stop := make(chan struct{})
	go handleInput(stop)
	Render(stop)

	fmt.Print("\033[H\033[2J\033[?25h")
}
