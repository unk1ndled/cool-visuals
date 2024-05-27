package tree

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	screenWidth  = 1200
	screenHeight = 800
)

var (
	factor       float64 = 0.623
	patternAngle float64 = 0
)

type Line struct {
	Startx, Starty, Endx, Endy float64
	angle                      float64
}

type RecursiveTree struct {
	branches    []*Line
	root        *Line
	initAngle   float64
	leavesCount int

	length   float64
	startLen float64
}

func NewRecursive(ang float64, len float64, rootx, rooty float64) *RecursiveTree {

	patternAngle = ang
	return &RecursiveTree{
		root:        &Line{Startx: rootx, Starty: rooty, Endx: rootx, Endy: rooty - len, angle: 0},
		branches:    make([]*Line, 0),
		initAngle:   patternAngle * math.Pi,
		leavesCount: 1,
		startLen:    len,
		length:      float64(len),
	}
}

func (r *RecursiveTree) Reset() {
	r.branches = make([]*Line, 0)
	r.leavesCount = 1
	r.initAngle = patternAngle * math.Pi
	r.length = r.startLen
}

func (r *RecursiveTree) ComputeChildren(line *Line, len float64) (*Line, *Line) {
	startx, starty := line.Endx, line.Endy

	//this took me way too much time to figure out
	//the angle is from the y axis to the x one
	//go right by incrementing the angle and left by decrementing it

	ra, la := line.angle+r.initAngle, line.angle-r.initAngle
	nextx, nexty := startx+(math.Sin(+ra)*len), starty-(math.Cos(+ra)*len)
	line1 := &Line{startx, starty, nextx, nexty, ra}
	nextx, nexty = startx+(math.Sin(la)*len), starty-(math.Cos(la)*len)
	line2 := &Line{startx, starty, nextx, nexty, la}
	return line1, line2
}

func (r *RecursiveTree) Calc(renderer *sdl.Renderer) {
	if len(r.branches) == 0 {
		l1, l2 := r.ComputeChildren(r.root, r.length)
		r.branches = append(r.branches, l1)
		r.branches = append(r.branches, l2)
	}
	if r.leavesCount >= 2 {
		leaves := r.branches[len(r.branches)-r.leavesCount:]
		for _, leave := range leaves {
			l1, l2 := r.ComputeChildren(leave, r.length)
			r.branches = append(r.branches, l1)
			r.branches = append(r.branches, l2)
		}
	}

	r.length *= factor
	r.leavesCount *= 2
}

func (r *RecursiveTree) Draw(renderer *sdl.Renderer) {
	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.DrawLine(int32(r.root.Startx), int32(r.root.Starty), int32(r.root.Endx), int32(r.root.Endy))
	for _, branche := range r.branches {
		renderer.DrawLine(int32(branche.Startx), int32(branche.Starty), int32(branche.Endx), int32(branche.Endy))
	}
}

func (r *RecursiveTree) Update(renderer *sdl.Renderer) bool {
	done := false
	if r.length > 1 {
		r.Calc(renderer)
	} else {
		done = true
	}
	r.Draw(renderer)
	return done

}
