package main

import (
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var people [256]*Person

var dt = 1. / 60.

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(people); i++ {
		people[i] = newPerson()
		people[i].Position = pixel.V(random(10, 1014), random(10, 758))
		people[i].Velocity = pixel.V(random(-50, 50), random(-50, 50))
	}

	imd := imdraw.New(nil)
	// imd.Precision = 7
	// imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	for !win.Closed() {
		imd.Clear()

		for _, p := range people {
			p.draw(imd)
			if p.Position.X-10 <= 0 || p.Position.X+10 >= 1024 {
				p.Velocity.X = -p.Velocity.X
			}
			if p.Position.Y-10 <= 0 || p.Position.Y+10 >= 768 {
				p.Velocity.Y = -p.Velocity.Y
			}
			p.update(dt)
		}

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func random(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	pixelgl.Run(run)
}
