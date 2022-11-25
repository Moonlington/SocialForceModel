package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var people [512]*Person
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

	for i := 0; i < len(people)/2; i++ {
		people[i] = newPerson(i)
		people[i].Position = pixel.V(random(-550, -450), random(-350, 350))
		people[i].DesiredSpeed = math.Max(0.5, (rand.NormFloat64()*0.25+1.33)*10)

		people[i].Goals = []*Goal{newGoal(pixel.V(0, 0), 100, 0), newGoal(pixel.V(random(400, 500), random(-200, 200)), 100, 0)}

		people[i].Radius = math.Max(.05, (rand.NormFloat64()*0.05+0.2)*25)
		people[i].Mass = math.Max(45, rand.NormFloat64()*5+70)
		people[i].wallThreshold = math.Max(people[i].Radius+0.05, (rand.NormFloat64()*0.1+2)*25)
	}

	for i := len(people) / 2; i < len(people); i++ {
		people[i] = newPerson(i)

		people[i].Color = colornames.Magenta

		people[i].Position = pixel.V(random(450, 550), random(-350, 350))
		people[i].DesiredSpeed = math.Max(0.5, (rand.NormFloat64()*0.25+1.33)*10)

		people[i].Goals = []*Goal{newGoal(pixel.V(0, 0), 100, 0), newGoal(pixel.V(random(-500, -400), random(-350, 350)), 100, 0)}

		people[i].Radius = math.Max(.05, (rand.NormFloat64()*0.05+0.2)*25)
		people[i].Mass = math.Max(45, rand.NormFloat64()*5+70)
		people[i].wallThreshold = math.Max(people[i].Radius+0.05, (rand.NormFloat64()*0.1+2)*25)
	}

	obstacles = append(obstacles, newObstacle(pixel.R(-50, 100, 50, 500), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-50, -500, 50, -100), false))

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

		for _, o := range obstacles {
			o.Draw(imd)
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
