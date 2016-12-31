package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

type Cell struct {
	CurrentState int
	NextState    int
}

type Coord struct {
	X int
	Y int
}

type GameOfLife struct {
	Width         int
	Height        int
	Cells         [][]Cell
	SelectedCoord Coord
}

func NewBoard(width int, height int) *GameOfLife {
	g := &GameOfLife{Width: width, Height: height}
	g.Cells = make([][]Cell, width)
	for i := 0; i < width; i++ {
		g.Cells[i] = make([]Cell, height)
		for j := 0; j < height; j++ {
			g.Cells[i][j].CurrentState = 0
			g.Cells[i][j].NextState = 0
		}
	}
	g.SelectedCoord = Coord{X: 0, Y: 0}
	return g
}

func (g *GameOfLife) Render() {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			draw := termbox.Cell{Ch: ' ', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack}
			if g.Cells[x][y].CurrentState == 1 {
				draw.Ch = 'x'
			}
			if x == g.SelectedCoord.X && y == g.SelectedCoord.Y {
				draw.Fg = termbox.ColorBlack
				draw.Bg = termbox.ColorWhite
			}
			termbox.SetCell(x, y, draw.Ch, draw.Fg, draw.Bg)
		}
	}
	str := strconv.Itoa(g.GetAliveCountAround(g.SelectedCoord.X, g.SelectedCoord.Y))
	s, _ := utf8.DecodeRuneInString(str[0:])
	termbox.SetCell(g.Width+10, 2, s, termbox.ColorWhite, termbox.ColorBlack)
}

func (g *GameOfLife) Update() {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			alive := g.GetAliveCountAround(x, y)
			if g.Cells[x][y].CurrentState == 1 {
				if alive < 2 || alive > 3 {
					g.Cells[x][y].NextState = 0
				} else {
					g.Cells[x][y].NextState = 1
				}
			} else {
				if alive == 3 {
					g.Cells[x][y].NextState = 1
				}
			}
		}
	}
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			g.Cells[x][y].CurrentState = g.Cells[x][y].NextState
		}
	}
}

func (g *GameOfLife) RandomizeBoard() {
	for x := 0; x < g.Width; x++ {
		for y := 0; y < g.Height; y++ {
			g.Cells[x][y].CurrentState = rand.Intn(2)
		}
	}
}

func (g *GameOfLife) GetAliveCountAround(cellX, cellY int) int {
	count := 0
	for x := cellX - 1; x < cellX+2; x++ {
		for y := cellY - 1; y < cellY+2; y++ {
			if x == cellX && y == cellY {
				continue
			}
			if x >= 0 && x < g.Width && y >= 0 && y < g.Height && g.Cells[x][y].CurrentState == 1 {
				count += 1
			}
		}
	}
	return count
}

func (g *GameOfLife) MoveCursorUp() {
	if g.SelectedCoord.Y > 0 {
		g.SelectedCoord.Y -= 1
	}
}
func (g *GameOfLife) MoveCursorDown() {
	if g.SelectedCoord.Y < g.Height-1 {
		g.SelectedCoord.Y += 1
	}
}
func (g *GameOfLife) MoveCursorLeft() {
	if g.SelectedCoord.X > 0 {
		g.SelectedCoord.X -= 1
	}
}
func (g *GameOfLife) MoveCursorRight() {
	if g.SelectedCoord.X < g.Width-1 {
		g.SelectedCoord.X += 1
	}
}

func (g *GameOfLife) SetSelectedCell(value int) {
	if g.SelectedCoord.X >= 0 && g.SelectedCoord.X < g.Width && g.SelectedCoord.Y >= 0 && g.SelectedCoord.Y < g.Height {
		g.Cells[g.SelectedCoord.X][g.SelectedCoord.Y].CurrentState = value
	}
}

func render(g *GameOfLife) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	g.Render()
	termbox.Flush()
}

func main() {

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	pause := false
	update := make(chan int)
	go func() {
		for {
			if !pause {
				update <- 1
			}
			time.Sleep(1 * time.Second / 3)
		}
	}()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()

	g := NewBoard(40, 20)
	g.RandomizeBoard()

loop:
	for {
		select {
		case e := <-events:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyEsc:
					break loop
				case termbox.KeyCtrlA:
					fmt.Println(e.Ch)
				case termbox.KeySpace:
					pause = !pause
				default:
					switch e.Ch {
					case 'h':
						g.MoveCursorLeft()
					case 'l':
						g.MoveCursorRight()
					case 'j':
						g.MoveCursorDown()
					case 'k':
						g.MoveCursorUp()
					case 'r':
						g.RandomizeBoard()
					case '1':
						g.SetSelectedCell(1)
					case '0':
						g.SetSelectedCell(0)
					}
				}
			}
			render(g)
		case _ = <-update:
			g.Update()
			render(g)
		}
	}

}
