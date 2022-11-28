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

var people [256]*Person
var obstacles []*Obstacle

func run() {
	rand.Seed(time.Now().Unix())
	cfg := pixelgl.WindowConfig{
		Title:  "Sociophysics Group 3 - Social Force Model",
		Bounds: pixel.R(0, 0, 1800, 800),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(people)/2; i++ {
		people[i] = newPerson(i)
		people[i].Position = pixel.V(random(-700, -800), random(-350, 350))
		people[i].Behavior = NewPathBehavior(newPath(newGoal(pixel.V(-350, 250), 100, 0), newGoal(pixel.V(random(700, 800), random(-350, 350)), 50, 0)))
	}

	for i := len(people) / 2; i < len(people); i++ {
		people[i] = newPerson(i)

		people[i].Color = colornames.Magenta
		people[i].Position = pixel.V(random(800, 700), random(-350, 350))
		people[i].Behavior = NewPathBehavior(newPath(newGoal(pixel.V(350, 250), 100, 0), newGoal(pixel.V(random(-800, -700), random(-350, 350)), 50, 0)))
	}

	// The last few from each group should follow the first person in their group
	amount := 10
	for i := 0; i < amount; i++ {
		people[i].Color = colornames.Darkcyan
		people[i].Behavior = NewFollowerBehavior(people[amount+1])
		people[i+len(people)/2].Color = colornames.Darkmagenta
		people[i+len(people)/2].Behavior = NewFollowerBehavior(people[len(people)/2+amount+1])
	}

	obstacles = append(obstacles, newObstacle(pixel.R(-300, -100, 300, 100), false)) // Kiosk

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
			p.Draw(win, imd)
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
