package sdlscreen

import (
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	ScreenWidth  int32 = 0
	ScreenHeight int32 = 0
)

type Runnable interface {
	Update(*sdl.Renderer) bool
}

func Visualise(name string, w, h int32, apps ...Runnable) {
	ScreenHeight = h
	ScreenWidth = w
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Failed to initialise SDL : %s\n", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Fprintf(os.Stderr, " Failed to Create window : %s\n", err)
		os.Exit(2)
	}
	defer window.Destroy()

	renderer, _ := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)

	quit := false

	for !quit {
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			if e.GetType() == sdl.QUIT {
				quit = true
			}
		}
		done := true
		for _, app := range apps {
			done = app.Update(renderer) && done
		}
		if done {
			sdl.Delay(3500)
			quit = true
			continue
		}

		//delete later
		//$$$$$$$$$$$$//
		msg := ""
		x, y := int32(50), int32(50)
		size := int32(5)
		spacing := size * 3
		for _, char := range msg {
			drawCharacter(renderer, char, x, y, size)
			if char == 'm' {
				x += spacing
			}
			if char == 'i' {
				x -= 2 * size
			}
			x += spacing
		}
		//$$$$$$$$$$$$//

		renderer.Present()
		sdl.Delay(0)
	}

}

// Define each character with line segments (pairs of coordinates)
var font = map[rune][][4]int32{
	'm': {{0, 4, 0, 0}, {0, 0, 2, 4}, {2, 4, 4, 0}, {4, 0, 4, 4}},
	'i': {{0, 0, 0, 4}},
	's': {{2, 0, 0, 0}, {0, 0, 0, 2}, {0, 2, 2, 2}, {2, 2, 2, 4}, {2, 4, 0, 4}},
	'y': {{0, 0, 2, 2}, {2, 2, 2, 4}, {2, 2, 4, 0}},
	'o': {{0, 0, 0, 4}, {0, 4, 2, 4}, {2, 4, 2, 0}, {2, 0, 0, 0}},
	'u': {{0, 0, 0, 4}, {0, 4, 2, 4}, {2, 4, 2, 0}},
	'a': {{0, 4, 2, 0}, {2, 0, 4, 4}},
}

func drawCharacter(renderer *sdl.Renderer, char rune, x, y, size int32) {
	lines, exists := font[char]
	if !exists {
		return
	}

	for _, line := range lines {
		renderer.DrawLine(
			x+line[0]*size, y+line[1]*size,
			x+line[2]*size, y+line[3]*size,
		)
	}
}
