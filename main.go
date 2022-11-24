package main

import (
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var people [256]*Person

func run() {
	rand.Seed(time.Now().Unix())
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
		people[i] = newPerson(i)
		people[i].Position = pixel.V(random(10, 250), random(10, 758))
		people[i].DesiredSpeed = random(10., 20.)
		people[i].Goal = pixel.V(random(750, 1014), random(10, 758))
	}

	imd := imdraw.New(nil)
	// imd.Precision = 7
	// imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

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
