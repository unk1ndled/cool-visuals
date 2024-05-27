package mb

import (
	"fmt"
	"sync"

	unkutil "github.com/unk1ndled/sdl-go/util"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	UP = iota
	DOWN
	LEFT
	RIGHT

	MaxIter = 500
	Bound   = 15
)

var (
	RangeIncement  = 0.2
	zoomMultiplier = 0.5
)

type Point struct {
	Real, Imag       float64
	screenX, screenY int32

	iter int
}

func (pt *Point) Add(other *Point) {
	pt.Real += other.Real
	pt.Imag += other.Imag
}
func (pt *Point) SquaredMagnitude() float64 {
	return pt.Imag*pt.Imag + pt.Real*pt.Real
}
func (pt *Point) Square(sContainer *Point) {
	sContainer.Real, sContainer.Imag = pt.Real*pt.Real-pt.Imag*pt.Imag, 2*pt.Real*pt.Imag
}

func (pt *Point) Compute() {
	acc := &Point{Real: 0, Imag: 0}
	var i int
	for i = 0; i < MaxIter; i++ {
		acc.Square(acc)
		acc.Add(pt)
		if acc.SquaredMagnitude() > Bound*Bound {
			break
		}
	}
	pt.iter = i
	if i == MaxIter {
		pt.iter = -1
	}
}

type Set struct {
	screenW  int
	screenH  int
	computed bool

	points []*Point

	xRangestart float64
	yRangestart float64

	xRangeend float64
	yRangeend float64
}

func NewSet(w, h int32) *Set {
	set := &Set{
		screenW: int(w),
		screenH: int(h),

		xRangestart: -2.0,
		yRangestart: -1.5,

		xRangeend: 1.0,
		yRangeend: 1.5,
		computed:  false,
	}
	set.points = make([]*Point, set.screenW*set.screenH)
	for i := 0; i < set.screenW; i++ {
		for j := 0; j < set.screenH; j++ {
			set.points[i*set.screenH+j] = &Point{screenX: int32(i), screenY: int32(j)}
		}
	}
	return set
}

func (set *Set) initialise() {
	var wg sync.WaitGroup
	numWorkers := 8
	workSize := set.screenW / numWorkers

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func(start, end int) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := 0; j < set.screenH; j++ {
					pt := set.points[i*set.screenH+j]
					pt.Real = unkutil.Map(float64(i), 0, float64(set.screenW), set.xRangestart, set.xRangeend)
					pt.Imag = unkutil.Map(float64(j), 0, float64(set.screenH), set.yRangestart, set.yRangeend)
					pt.Compute()
				}
			}
		}(w*workSize, (w+1)*workSize)
	}

	wg.Wait()
	set.computed = true
}

type Color struct {
	R, G, B uint8
}

func createPalette(size int) []Color {
	palette := make([]Color, size)
	for i := 0; i < size; i++ {
		t := float64(i) / float64(size)
		palette[i] = Color{
			R: uint8(9 * (1 - t) * t * t * t * 255),
			G: uint8(15 * (1 - t) * (1 - t) * t * t * 255),
			B: uint8(8.5 * (1 - t) * (1 - t) * (1 - t) * t * 255),
		}
	}
	return palette
}

func (set *Set) Draw(rdr *sdl.Renderer) {
	palette := createPalette(MaxIter)
	for _, pt := range set.points {
		if pt.iter == -1 {
			rdr.SetDrawColor(0, 0, 0, 255) // Black for points inside the set
		} else {
			color := palette[pt.iter%MaxIter]
			rdr.SetDrawColor(color.R, color.G, color.B, 255)
		}
		rdr.DrawPoint(pt.screenX, pt.screenY)
	}
}

func (set *Set) Zoom(in bool) {
	multiplier := zoomMultiplier
	if !in {
		multiplier = 1 / zoomMultiplier
	}

	xRangeMid := (set.xRangestart + set.xRangeend) / 2
	yRangeMid := (set.yRangestart + set.yRangeend) / 2

	newWidth := (set.xRangeend - set.xRangestart) * multiplier
	newHeight := (set.yRangeend - set.yRangestart) * multiplier

	set.xRangestart = xRangeMid - newWidth/2
	set.xRangeend = xRangeMid + newWidth/2
	set.yRangestart = yRangeMid - newHeight/2
	set.yRangeend = yRangeMid + newHeight/2

	RangeIncement = newWidth / 10
}

func (set *Set) Translate(axis int) {
	switch axis {
	case UP:
		set.yRangestart -= RangeIncement
		set.yRangeend -= RangeIncement

	case RIGHT:
		set.xRangestart += RangeIncement
		set.xRangeend += RangeIncement

	case LEFT:
		set.xRangestart -= RangeIncement
		set.xRangeend -= RangeIncement
	case DOWN:
		set.yRangestart += RangeIncement
		set.yRangeend += RangeIncement
	}
}

func (set *Set) HandleInput() bool {
	keys := sdl.GetKeyboardState()
	updated := false
	if keys[sdl.SCANCODE_UP] != 0 {
		set.Translate(UP)
		updated = true
	} else if keys[sdl.SCANCODE_DOWN] != 0 {
		set.Translate(DOWN)
		updated = true
	} else if keys[sdl.SCANCODE_LEFT] != 0 {
		set.Translate(LEFT)
		updated = true
	} else if keys[sdl.SCANCODE_RIGHT] != 0 {
		set.Translate(RIGHT)
		updated = true
	} else if keys[sdl.SCANCODE_Q] != 0 {
		set.Zoom(false)
		updated = true
	} else if keys[sdl.SCANCODE_W] != 0 {
		set.Zoom(true)
		updated = true
	}
	return updated
}

func (set *Set) Update(rdr *sdl.Renderer) bool {
	updated := set.HandleInput()

	if updated {
		set.computed = false
	}

	if !set.computed {
		set.initialise()
		fmt.Println("init")
	}
	set.Draw(rdr)

	return false
}
