package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var people [500]*Person
var obstacles []*Obstacle

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
		people[i].Position = pixel.V(random(-500, -400), random(-350, 350))
		people[i].DesiredSpeed = (rand.NormFloat64()*0.1 + 1.33) * 10
		people[i].Goal = pixel.V(random(400, 500), random(-200, 200))
		people[i].Radius = (rand.NormFloat64()*0.05 + 0.2) * 25
		people[i].Mass = rand.NormFloat64()*5 + 70

		people[i].wallThreshold = (rand.NormFloat64()*0.1 + 1) * 25
	}

	obstacles = append(obstacles, newObstacle(pixel.R(-50, -100, 50, 100), false))

	imd := imdraw.New(nil)
	// imd.Precision = 7
	imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	last := time.Now()

	for !win.Closed() {
		imd.Clear()

		dt := time.Since(last).Seconds()
		last = time.Now()

		wg := new(sync.WaitGroup)
		for _, p := range people {
			p.Draw(imd)
			wg.Add(1)
			go func(p *Person) {
				defer wg.Done()
				p.update(dt, people[:], obstacles[:])
			}(p)
		}
		wg.Wait()

		// for _, o := range obstacles {
		// 	o.Draw(imd)
		// }

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
