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

var people [0]*Person
var obstacles []*Obstacle
var edges []*Obstacle
var emptybins *EmptyBin[*Person] = newEmptyBin[*Person](10, 5, -900, 900, -400, 400)

var triangulation *Triangulation

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

	createObstaclesAndEdges()

	wanderLocations := []*Goal{}
	wanderLocations = generateWanderLocations(wanderLocations)

	// Using the list of points from wanderLocations, create a triangulation
	points := []pixel.Vec{}
	for _, goal := range wanderLocations {
		points = append(points, goal.Target)
	}
	triangulation = BowyerWatson(points)
	// fmt.Println(triangulation)

	// people[i].Behavior = NewPathBehavior(newPath(newGoal(pixel.V(-350, 250), 100, 0), newGoal(pixel.V(random(700, 800), random(-350, 350)), 50, 0)))
	// people[i].Behavior = NewPathBehavior(newPath(newGoal(pixel.V(350, 250), 100, 0), newGoal(pixel.V(random(-800, -700), random(-350, 350)), 50, 0)))
	// The last few from each group should follow the first person in their group
	createPeople(wanderLocations)

	for _, person := range people {
		emptybins.Add(person)
	}

	imd := imdraw.New(nil)
	imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	last := time.Now()

	for !win.Closed() {
		imd.Clear()

		dt := time.Since(last).Seconds()
		last = time.Now()

		// p.update(dt, people[:], obstacles[:])
		updateAndDrawPeople(win, imd, dt)

		emptybins.Update()

		for _, o := range obstacles {
			o.Draw(imd)
		}

		triangulation.Draw(imd)

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
}

func updateAndDrawPeople(win *pixelgl.Window, imd *imdraw.IMDraw, dt float64) {
	wg := new(sync.WaitGroup)
	for _, p := range people {
		p.Draw(win, imd)
		wg.Add(1)
		go func(p *Person) {
			defer wg.Done()

			p.update(dt, emptybins.GetSurrounding(p, 1), obstacles[:])
		}(p)
	}
	wg.Wait()
}

func createPeople(wanderLocations []*Goal) {
	for i := 0; i < len(people)/2; i++ {
		people[i] = newPerson(i)
		people[i].Position = pixel.V(random(-400, -800), random(-150, 150))
		noCollision := true
		for noCollision {
			noCollision = false
			for j := 0; j < i; j++ {
				if people[i].Position.To(people[j].Position).Len() < people[i].Radius+people[j].Radius*1.1 {
					people[i].Position = pixel.V(random(-400, -800), random(-150, 150))
					noCollision = true
					break
				}
			}
		}

		people[i].Behavior = NewWanderBehavior(obstacles, wanderLocations...)
	}

	for i := len(people) / 2; i < len(people); i++ {
		people[i] = newPerson(i)

		people[i].Color = colornames.Magenta
		people[i].Position = pixel.V(random(800, 400), random(-150, 150))
		noCollision := true
		for noCollision {
			noCollision = false
			for j := len(people) / 2; j < i; j++ {
				if people[i].Position.To(people[j].Position).Len() < people[i].Radius+people[j].Radius*1.1 {
					people[i].Position = pixel.V(random(800, 400), random(-150, 150))
					noCollision = true
					break
				}
			}
		}

		people[i].Behavior = NewWanderBehavior(obstacles, wanderLocations...)
	}

	if len(people) > 0 {
		amount := 4
		for i := 0; i < amount; i++ {
			people[i].Color = colornames.Darkcyan
			people[i].Behavior = NewFollowerBehavior(people[amount+1], obstacles)
			people[i+len(people)/2].Color = colornames.Darkmagenta
			people[i+len(people)/2].Behavior = NewFollowerBehavior(people[len(people)/2+amount+1], obstacles)
		}
	}
}

func generateWanderLocations(wanderLocations []*Goal) []*Goal {
	for i := 0; i < 200; i++ {
		wanderLocations = append(wanderLocations, NewGoal(pixel.V(random(400, 800), random(-180, 180)), 250, rand.NormFloat64()*0.5+10))
		wanderLocations = append(wanderLocations, NewGoal(pixel.V(random(-800, -400), random(-180, 180)), 250, rand.NormFloat64()*0.5+10))
	}
	// wanderLocations = append(wanderLocations, NewGoal(pixel.V(-200, 250), 100, 0), NewGoal(pixel.V(200, 250), 100, 0))
	return wanderLocations
}

func createObstaclesAndEdges() {
	obstacles = append(obstacles, newObstacle(pixel.R(-890, 200, 890, 390), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-890, -390, 890, -200), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-150, -100, 150, 100), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-890, -390, 890, 390), true))

	edges = append(edges, newObstacle(pixel.R(-890, 200, 890, 390), false))
	edges = append(edges, newObstacle(pixel.R(-890, -390, 890, -200), false))
}

func random(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	pixelgl.Run(run)
}
