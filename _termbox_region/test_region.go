package main

import (
	"time"

	termbox "github.com/nsf/termbox-go"
	tbregion "github.com/zwodahs/termbox-go-region"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	tbregion.InitBorder()
	defer termbox.Close()

	events := make(chan termbox.Event)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

	update := make(chan int)
	go func() {
		for {
			update <- 1
			time.Sleep(1 * time.Second / 3)
		}
	}()

	termbox.SetInputMode(termbox.InputEsc)
	termbox.Flush()

	mainRegion := tbregion.NewRegion(100, 100, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorDefault})

	region := mainRegion.NewRegion(10, 10, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorDefault})
	region.DrawThinBorder()

	region2 := region.NewRegion(5, 5, termbox.Cell{Ch: ' ', Fg: termbox.ColorDefault, Bg: termbox.ColorDefault})
	region2.DrawThinBorder()
loop:
	for {
		select {
		case e := <-events:
			switch e.Type {
			case termbox.EventKey:
				switch e.Key {
				case termbox.KeyEsc:
					break loop
				default:
					switch e.Ch {
					case '1':
						position := region.GetPosition()
						position[0] += 1
						region.SetPosition(position)
					case '2':
						position := region2.GetPosition()
						position[0] += 1
						region2.SetPosition(position)
					case 'h':
						region2.Hidden = !region2.Hidden
					}
				}
			}
			mainRegion.Draw()
			termbox.Flush()
		case <-update:
			mainRegion.Draw()
			termbox.Flush()
		}
	}
}
