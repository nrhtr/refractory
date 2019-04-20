package main

import (
	//"github.com/nsf/termbox-go"
	"fmt"
	"log"
	"os/exec"
)

type ItemKind int8
type Item struct {
	kind ItemKind
	//flags int8
	count int16
}

type Widget int8

const (
	EmptyI ItemKind = iota
	Charcoal
	Hematite
	Magnetite
	Log
	Limestone
	Chalk
	Marble
	Granite
)

const (
	EmptyW Widget = iota
	Wall
	Furnace
	BeltU
	BeltD
	BeltL
	BeltR
)


type gameState struct {
	width     int32
	height    int32
	px        int32 // assume single player
	py        int32
	items     []Item
	widgets   []Widget
}

func defaultgameState(w int32, h int32) gameState {
	g := gameState{
		width:     w,
		height:    h,
		px:        0,
		py:        0,
		items:     make([]Item, w*h),
		widgets:   make([]Widget, w*h),
	}

	for i := range g.widgets {
		g.widgets[i] = EmptyW
	}

	for x := int32(0); x < g.width; x++ {
		placeWidget(g, x, 0,          Wall)
		placeWidget(g, x, g.height-1, Wall)
	}

	for y := int32(0); y < g.height; y++ {
		placeWidget(g, 0,         y, Wall)
		placeWidget(g, g.width-1, y, Wall)
	}

	return g
}

func printGameBoard(g gameState) {
	for y := int32(0); y < g.height; y++ {
		for x := int32(0); x < g.width; x++ {
			i := g.items[y * g.width + x]
			w := g.widgets[y * g.width + x]
			var c rune
			switch i.kind {
				case EmptyI: c = ' '
				case Log:    c = 'w'
				default:     c = '?'
			}
			if w != EmptyW {
				switch w {
				case EmptyW: c = ' '
				case Furnace:c = 'F'
				case BeltU:  c = '^'
				case BeltD:  c = '~'
				case BeltL:  c = '<'
				case BeltR:  c = '>'
				case Wall:   c = '#'
				default:     c = '?'
				}
			}
			fmt.Printf("%c", c)
		}
		fmt.Println()
	}
}

func placeItem(g gameState, x int32, y int32, i Item) {
	g.items[y*g.width+x] = i
}

func placeWidget(g gameState, x int32, y int32, w Widget) {
	g.widgets[y*g.width+x] = w
}

func stringToWidget(s string) []Widget {
	l := len(s)
	r := make([]Widget, l)
	for n := range s {
		var w Widget
		switch s[n] {
			// TODO: have single map between Widget <--> representation
			case ' ': w = EmptyW
			case 'F': w = Furnace
			case '^': w = BeltU
			case '~': w = BeltD
			case '<': w = BeltL
			case '>': w = BeltR
			case '#': w = Wall
		}
		r[n] = w
	}

	return r
}

func paintWidgetH(g gameState, x int32, y int32, m string) {
	start := x
	end   := start + int32(len(m))

	if (start < 0) || (end >= g.width) || (y < 0) || (y >= g.height) {
		log.Fatal("Widget collides with game boundary")
	}

	w := stringToWidget(m)

	for i := range w {
		placeWidget(g, start + int32(i), y, w[i])
	}
}

func paintWidgetV(g gameState, x int32, y int32, m string) {
	start := y
	end   := start + int32(len(m))

	if (start < 0) || (end >= g.height) || (x < 0) || (x >= g.width) {
		log.Fatal("Widget collides with game boundary")
	}

	w := stringToWidget(m)

	for i := range w {
		placeWidget(g, x, start + int32(i), w[i])
	}
}

func stepGame(g gameState) gameState {
	newState := g
	for i := range g.widgets {
		x := i % g.width
		y := i / g.height
		switch g.widgets[i] {
		case BeltR:
			speed := 5
			if g.items[i] != EmptyI {
				tryMoveItems(newState, i, speed
			}
		case BeltL: newState.widgets[i] = BeltR
		default: newState.widgets[i] = g.widgets[i]
		}
	}

	return newState
}

func main() {
	//game := gameState{width: 32, height: 32, px: 0, py: 0, items:
	game := defaultgameState(64, 32)
	fmt.Println("Created new game of", game.width, "by", game.height)
	fmt.Println("The size of the map is", len(game.items), "cells")
	placeItem(game, 30, 20, Item{Log, 16})
	paintWidgetH(game, 5, 10, ">>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	paintWidgetV(game, 5, 11, "^^^^^^^^^^")
	paintWidgetH(game, 5, 21, "<<<<<<<<<<")

	printGameBoard(game)
	game = stepGame(game)
	exec.Command("clear")

	printGameBoard(game)
	game = stepGame(game)
	exec.Command("clear")

	printGameBoard(game)
	game = stepGame(game)
	exec.Command("clear")
}
