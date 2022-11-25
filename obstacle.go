package main

import (
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Obstacle struct {
	pixel.Rect

	Inner bool
}

func newObstacle(r pixel.Rect, inner bool) *Obstacle {
	o := new(Obstacle)

	o.Rect = r
	o.Inner = inner

	return o
}

func (o *Obstacle) Dist(p *Person) pixel.Vec {
	shortestVec := pixel.V(math.Inf(1), math.Inf(1))
	for _, e := range o.Edges() {
		distVec := p.Position.To(e.Closest(p.Position))
		if shortestVec.Len() > distVec.Len() {
			shortestVec = distVec
		}
	}
	if !o.Inner && o.Contains(p.Position) {
		shortestVec = shortestVec.Scaled(-1)
	}
	return shortestVec
}

func (o *Obstacle) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(o.Min)
	imd.Push(o.Max)
	imd.Rectangle(1)
}
