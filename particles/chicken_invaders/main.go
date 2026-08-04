package main

import (
	"fmt"
	"golang.org/x/term"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	w = 25
	h = 15
)

var mu sync.Mutex
var jett Jett
var screen [][]rune
var bullets []*Bullet
var chickens []*Chicken
var score int32 = 0

type Jett struct {
	x, y float32
	char rune
}

type Bullet struct {
	x, y float32
	char rune
}

type Chicken struct {
	x, y float32
	char rune
}

func (c *Chicken) goDown() {
	c.y += 0.5
}

func (b *Bullet) goUp() {
	b.y -= 1
}

func (j *Jett) emitBullet() {
	bullet := CreateBullet(j.x, j.y-1)
	bullets = append(bullets, bullet)
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

func CreateBullet(x, y float32) *Bullet {
	return &Bullet{
		x:    x,
		y:    y,
		char: '|',
	}
}

func CreateChicken(x, y float32) *Chicken {
	return &Chicken{
		x:    x,
		y:    y,
		char: 'A',
	}
}

func launchChickens() {
	for {
		idx := rand.Intn(w)
		chicken := CreateChicken(float32(idx), 1)
		mu.Lock()
		chickens = append(chickens, chicken)
		mu.Unlock()

		time.Sleep(1 * time.Second)
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

func UpdateScreen(stop chan struct{}) {
	mu.Lock()
	defer mu.Unlock()

	// Render jett on screen
	if int(jett.x) >= 1 && int(jett.x) <= w-1 && int(jett.y) >= 1 && int(jett.y) <= h-1 {
		screen[int(jett.y)][int(jett.x)] = jett.char
	}

	// Check jett and chicken collision and stop the game
	for _, chicken := range chickens {
		if int(chicken.x) == int(jett.x) && int(chicken.y) == int(jett.y) {
			close(stop)
		}
	}

	bulletHits := make(map[int]bool)
	chickenHits := make(map[int]bool)

	checkColiisions(bulletHits, chickenHits)

	updateBullets(bulletHits)
	updateChickens(chickenHits)

	// Render bullets on screen and move them one step up
	for _, bullet := range bullets {
		if int(bullet.x) >= 1 && int(bullet.x) <= w-1 && int(bullet.y) >= 1 && int(bullet.y) <= h-1 {
			screen[int(bullet.y)][int(bullet.x)] = bullet.char
			bullet.goUp()
		}
	}

	// Render chickens on screen and move them one step down
	for _, chicken := range chickens {
		if int(chicken.x) >= 1 && int(chicken.x) <= w-1 && int(chicken.y) >= 1 && int(chicken.y) <= h-1 {
			screen[int(chicken.y)][int(chicken.x)] = chicken.char
			chicken.goDown()
		}
	}
}

func checkColiisions(bulletHits, chickenHits map[int]bool) {
	// If a bullet and a chicken share the same coordinates, they hit each other
	// remove both of them from the environment
	for i, bullet := range bullets {
		for j, chicken := range chickens {
			if int(bullet.x) == int(chicken.x) && int(bullet.y) <= int(chicken.y) {
				bulletHits[i] = true
				chickenHits[j] = true
				score += 1
				break
			}
		}
	}
}

func updateBullets(bulletHits map[int]bool) {
	// If a bullet was hit, or went outside the coordinates,
	// remove it from the environment
	var newBullets []*Bullet
	for i, bullet := range bullets {
		if !bulletHits[i] && int(bullet.y) > 0 {
			newBullets = append(newBullets, bullet)
		}
	}
	bullets = newBullets
}

func updateChickens(chickenHits map[int]bool) {
	// If a chicken was hit, or went outside the coordinates,
	// remove it from the environment
	var newChickens []*Chicken
	for i, chicken := range chickens {
		if !chickenHits[i] && int(chicken.y) < h-1 {
			newChickens = append(newChickens, chicken)
		}
	}
	chickens = newChickens
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
				mu.Lock()
				jett.goUp()
				mu.Unlock()
			case 'B':
				mu.Lock()
				jett.goDown()
				mu.Unlock()
			case 'C':
				mu.Lock()
				jett.goRight()
				mu.Unlock()
			case 'D':

				mu.Lock()
				jett.goLeft()
				mu.Unlock()
			}
		} else if buf[0] == 32 {
			mu.Lock()
			jett.emitBullet()
			mu.Unlock()
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
		UpdateScreen(stop)

		fmt.Print("\033[H\033[2J")
		for _, row := range screen {
			fmt.Print(string(row) + "\r\n")
		}
		fmt.Println(score)

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
	go launchChickens()

	stop := make(chan struct{})
	go handleInput(stop)
	Render(stop)

	fmt.Print("\033[H\033[2J\033[?25h")
}
