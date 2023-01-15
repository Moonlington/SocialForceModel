package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var peopleAmount int
var outputName string

const maxTimeSpend time.Duration = time.Minute * 5
const nudge = true

func init() {
	flag.IntVar(&peopleAmount, "a", 64, "Amount of people")
	flag.StringVar(&outputName, "o", "data.csv", "Output for the file")
}

var people []*Person
var obstacles []*Obstacle
var edges []*Obstacle
var emptybins *EmptyBin[*Person] = newEmptyBin[*Person](10, 5, -900, 900, -400, 400)

var secondsFromStart float64

var triangulation *Triangulation

func run() {
	flag.Parse()
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

	fmt.Println("Creating obstacles")
	createObstaclesAndEdges()

	fmt.Println("Generating wander locations")
	wanderLocations := generateWanderLocations()

	// Using the list of points from wanderLocations, create a triangulation
	fmt.Println("Generating triangulation")
	triangulation = BowyerWatson(wanderLocations)
	// fmt.Println(triangulation)
	// The last few from each group should follow the first person in their group
	fmt.Println("Generating people")
	createPeople()

	fmt.Println("Generating emptybin")
	for _, person := range people {
		emptybins.Add(person)
	}

	imd := imdraw.New(nil)
	imd.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

	// last := time.Now()

	for !win.Closed() && secondsFromStart <= maxTimeSpend.Seconds() {
		imd.Clear()

		// dt := time.Since(last).Seconds()
		dt := 3 * time.Second.Seconds() / 60
		// last = time.Now()
		secondsFromStart += dt

		// p.update(dt, people[:], obstacles[:])
		updateAndDrawPeople(win, imd, dt)

		emptybins.Update()

		for _, o := range obstacles {
			o.Draw(imd)
		}

		// triangulation.Draw(imd)

		writeToData()

		win.Clear(colornames.Black)
		imd.Draw(win)
		win.Update()
	}
	file, err := os.Create(outputName)
	if err != nil {
		panic("wtf file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.WriteAll(data)
}

var data [][]string

func writeToData() {
	for _, person := range people {
		personData := []string{}
		personData = append(personData, fmt.Sprintf("%d", person.id), fmt.Sprintf("%f", secondsFromStart), fmt.Sprintf("%f", person.Position.X), fmt.Sprintf("%f", person.Position.Y), fmt.Sprintf("%t", person.Position.Y >= 170 || person.Position.Y <= -170))
		data = append(data, personData)
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

func createPeople() {
	for i := 0; i < peopleAmount/2; i++ {
		people = append(people, newPerson(i))
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

		people[i].Behavior = NewPathfinderBehavior(triangulation, obstacles)
	}

	for i := peopleAmount / 2; i < peopleAmount; i++ {
		people = append(people, newPerson(i))

		people[i].Color = colornames.Magenta
		people[i].Position = pixel.V(random(800, 400), random(-150, 150))
		noCollision := true
		for noCollision {
			noCollision = false
			for j := peopleAmount / 2; j < i; j++ {
				if people[i].Position.To(people[j].Position).Len() < people[i].Radius+people[j].Radius*1.1 {
					people[i].Position = pixel.V(random(800, 400), random(-150, 150))
					noCollision = true
					break
				}
			}
		}

		people[i].Behavior = NewPathfinderBehavior(triangulation, obstacles)
	}

	amount := peopleAmount / 16
	if peopleAmount > 2*amount+2 {
		for i := 0; i < amount; i++ {
			people[i].Color = colornames.Darkcyan
			people[i].Behavior = NewFollowerBehavior(people[amount+1], obstacles)
			people[i+peopleAmount/2].Color = colornames.Darkmagenta
			people[i+peopleAmount/2].Behavior = NewFollowerBehavior(people[peopleAmount/2+amount+1], obstacles)
		}
	}
}

func generateWanderLocations() []pixel.Vec {
	var wanderLocations []pixel.Vec
	if nudge {
		for i := 0; i < 50; i++ {
			wanderLocations = append(wanderLocations, pixel.V(random(400, 800), random(-150, 150)))
			wanderLocations = append(wanderLocations, pixel.V(random(-800, -400), random(-150, 150)))
		}
		wanderLocations = append(wanderLocations, pixel.V(-200, 170), pixel.V(200, 170))
		wanderLocations = append(wanderLocations, pixel.V(-200, -170), pixel.V(200, -170))
	} else {
		for i := 0; i < 400; i++ {
			test := pixel.V(random(-800, 800), random(-150, 150))
			collides := false
			for _, obstacle := range obstacles {
				if obstacle.Contains(test) && !obstacle.Inner {
					collides = true
					break
				}
			}
			if collides {
				continue
			}
			wanderLocations = append(wanderLocations, test)
		}
	}
	return wanderLocations
}

func createObstaclesAndEdges() {
	obstacles = append(obstacles, newObstacle(pixel.R(-890, 200, 890, 390), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-890, -390, 890, -200), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-150, -100, 150, 100), false))
	obstacles = append(obstacles, newObstacle(pixel.R(-890, -390, 890, 390), true))
	// obstacles = append(obstacles, newObstacle(pixel.R(-700, -5, -300, 5), false))

	edges = append(edges, newObstacle(pixel.R(-890, 200, 890, 390), false))
	edges = append(edges, newObstacle(pixel.R(-890, -390, 890, -200), false))
}

func random(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func main() {
	pixelgl.Run(run)
}
