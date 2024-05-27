package main

import (
	"github.com/unk1ndled/sdl-go/mb"
	sdlscreen "github.com/unk1ndled/sdl-go/sdl-screen"
)

func main() {
	w := 700.0
	h := 700.0
	/*

		for {
			var num1, num2 int

			// Move cursor to the top-left corner
			fmt.Print("\033[2J")
			fmt.Print("\033[H")
			fmt.Println("Enter two numbers:")

			// Read the two numbers from standard input
			_, err := fmt.Scan(&num1, &num2)
			if err != nil {
				fmt.Println("Error reading input:", err)
				return
			}
			result := float64(num1) / float64(num2)
			trees := make([]sdlscreen.Runnable, 1)

			trees[0] = tree.NewRecursive(result, 150, w/2, h-100)
			// trees[1] = tree.NewRecursive(result, 70, (w/2)-150, h)
			// trees[2] = tree.NewRecursive(result, 70, (w/2)+150, h)
			// trees[3] = tree.NewRecursive(result, 15, (w/2)-300, h)
			// trees[4] = tree.NewRecursive(result, 15, (w/2)+300, h)

			sdlscreen.Visualise("TREEE", int32(w), int32(h), trees...)
		}
	*/

	sdlscreen.Visualise("Mandle", int32(w), int32(h), mb.NewSet(int32(w), int32(h)))

}
