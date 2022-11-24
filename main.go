package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var people [200]*Person

func run() {
	rand.Seed(time.Now().Unix())
	cfg := pixelgl.WindowConfig{
		Title:  "Sociophysics Group 3 - Social Force Model",
		Bounds: pixel.R(0, 0, 1200, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(people); i++ {
		people[i] = newPerson(i)
		people[i].Position = pixel.V(random(-600, -500), random(-400, 400))
		people[i].DesiredSpeed = (rand.NormFloat64()*0.25 + 1.33) * 10
		people[i].Goal = pixel.V(random(500, 600), random(-400, 400))
		people[i].Radius = (rand.NormFloat64()*0.05 + 0.2) * 10
		people[i].Mass = rand.NormFloat64()*10 + 70
	}

	imd := imdraw.New(nil)
	// imd.Precision = 7
	imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	last := time.Now()

	for !win.Closed() {
		imd.Clear()

		dt := time.Since(last).Seconds()
		last = time.Now()

		for _, p := range people {
			p.draw(imd)
			// p.checkBoundaries(win.Bounds().Max)
			p.update(dt, people[:])
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
