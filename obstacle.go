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
	wallCircle := pixel.C(p.Position, p.wallThreshold)
	// fmt.Println(o, p)
	dist := o.IntersectCircle(wallCircle)
	if o.Inner {
		return dist.Unit().Scaled(dist.Len() - p.wallThreshold)
	} else {
		if dist.Len() == 0 {
			return dist.Unit().Scaled(math.Inf(1))
		}
		return dist.Unit().Scaled(p.wallThreshold - dist.Len())
	}
}

func (o *Obstacle) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(o.Min)
	imd.Push(o.Max)
	imd.Rectangle(1)
}
