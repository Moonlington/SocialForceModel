package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"golang.org/x/exp/slices"
)

// Behavior defines the behavior of a person.
type Behavior interface {
	GetTarget(p *Person, dt float64) pixel.Vec
}

// GoalBehavior defines the behavior of a person that goes to a goal.
type GoalBehavior struct {
	goal           *Goal
	maxRangeFactor float64
	closeEnough    bool
	LoiterTime     float64
}

func NewGoalBehavior(goal *Goal) *GoalBehavior {
	return &GoalBehavior{
		goal:           goal,
		maxRangeFactor: 1.5,
		closeEnough:    false,
	}
}

// SetMaxRangeFactor sets the max range factor.
func (b *GoalBehavior) SetMaxRangeFactor(maxRangeFactor float64) {
	b.maxRangeFactor = maxRangeFactor
}

// SetGoal sets the goal of the behavior.
func (b *GoalBehavior) SetGoal(goal *Goal) {
	b.goal = goal
}

// GetTarget gets the target of the behavior.
func (b *GoalBehavior) GetTarget(p *Person, dt float64) pixel.Vec {
	if b.goal == nil {
		return p.Position
	}
	if p.Position.To(b.goal.Target).Len() <= b.goal.Range {
		b.closeEnough = true
	}
	if p.Position.To(b.goal.Target).Len() <= b.goal.Range*b.maxRangeFactor && b.closeEnough {
		b.LoiterTime += dt
		return p.Position
	}
	b.closeEnough = false
	return b.goal.Target
}

// HasLoitered returns true if the person has loitered for the current goal.
func (b *GoalBehavior) HasLoitered() bool {
	return b.LoiterTime > b.goal.LoiterAfter
}

// Arrived returns true if the person has arrived at the goal.
func (b *GoalBehavior) Arrived() bool {
	return b.closeEnough
}

// FollowerBehavior defines the behavior of a person that follows a person.
type FollowerBehavior struct {
	Target       *Person
	Obstacles    []*Obstacle
	lastSeen     pixel.Vec
	goalBehavior *GoalBehavior
}

// NewFollowerBehavior creates a new follower behavior.
func NewFollowerBehavior(target *Person, obstacles []*Obstacle) *FollowerBehavior {
	goalB := NewGoalBehavior(NewGoal(target.Position, 0, 0))
	goalB.SetMaxRangeFactor(2)
	return &FollowerBehavior{
		Target:       target,
		Obstacles:    obstacles,
		goalBehavior: goalB,
	}
}

func (b *FollowerBehavior) SetTarget(target *Person) {
	b.Target = target
}

func (b *FollowerBehavior) GetTarget(p *Person, dt float64) pixel.Vec {
	if b.Target == nil {
		return p.Position
	}
	intersects := false
	for _, obstacle := range b.Obstacles {
		if obstacle.Inner {
			continue
		}
		if obstacle.IntersectLine(pixel.L(p.Position, b.Target.Position)).Len() > 0 {
			intersects = true
			break
		}
	}
	if !intersects {
		b.lastSeen = b.Target.Position
		b.goalBehavior.SetGoal(NewGoal(b.lastSeen, 1.5*(p.Radius+b.Target.Radius), 0))
	}
	return b.goalBehavior.GetTarget(p, dt)
}

// WanderBehavior defines the behavior of a person that walks to a random goal in sight
type WanderBehavior struct {
	WanderGoals  []*Goal
	CurrentGoal  *Goal
	Obstacles    []*Obstacle
	goalBehavior *GoalBehavior
}

// NewWanderBehavior creates a new wander behavior.
func NewWanderBehavior(obstacles []*Obstacle, wanderLocations ...*Goal) *WanderBehavior {
	return &WanderBehavior{WanderGoals: wanderLocations, Obstacles: obstacles, goalBehavior: NewGoalBehavior(nil)}
}

// Update updates the behavior.
func (b *WanderBehavior) GetTarget(p *Person, dt float64) pixel.Vec {
	if b.CurrentGoal == nil || b.goalBehavior.HasLoitered() {
		b.goalBehavior.LoiterTime = 0
		b.CurrentGoal = b.ChooseNextWanderLocation(p)
		b.goalBehavior.SetGoal(b.CurrentGoal)
	}
	return b.goalBehavior.GetTarget(p, dt)
}

func lineCollidesObstacles(A, B pixel.Vec, obstacles []*Obstacle) bool {
	for _, obstacle := range obstacles {
		if obstacle.Inner {
			continue
		}
		if obstacle.IntersectLine(pixel.L(A, B)).Len() > 0 {
			return true
		}
	}
	return false
}

func pointInObstacle(p pixel.Vec, obstacles []*Obstacle) bool {
	for _, obstacle := range obstacles {
		if obstacle.Inner {
			continue
		}
		if obstacle.Contains(p) {
			return true
		}
	}
	return false
}

// ChooseNextWanderLocation chooses the next wander location.
func (b *WanderBehavior) ChooseNextWanderLocation(p *Person) *Goal {
	var possibleGoals []*Goal
	for _, goal := range b.WanderGoals {
		if !lineCollidesObstacles(p.Position, goal.Target, b.Obstacles) {
			possibleGoals = append(possibleGoals, goal)
		}
	}
	if len(possibleGoals) == 0 {
		return nil
	}
	return possibleGoals[rand.Intn(len(possibleGoals))]
}

// PathBehavior defines the behavior of a person that follows a path.
type PathBehavior struct {
	Path         *Path
	CurrentGoal  *Goal
	GoalBehavior *GoalBehavior
}

// NewPathBehavior creates a new path behavior.
func NewPathBehavior(path *Path) *PathBehavior {
	return &PathBehavior{Path: path, GoalBehavior: NewGoalBehavior(nil)}
}

// GetTarget gets the target of the behavior.
func (b *PathBehavior) GetTarget(p *Person, dt float64) pixel.Vec {
	if b.CurrentGoal == nil || b.GoalBehavior.HasLoitered() {
		b.GoalBehavior.LoiterTime = 0
		if b.Path.Empty() {
			return b.GoalBehavior.GetTarget(p, dt)
		}
		b.CurrentGoal = b.Path.GetNextGoal()
		b.GoalBehavior.SetGoal(b.CurrentGoal)
	}
	return b.GoalBehavior.GetTarget(p, dt)
}

// SetPath sets the path of the behavior.
func (b *PathBehavior) SetPath(path *Path) {
	b.Path = path
	b.CurrentGoal = nil
}

// PathfinderBehavior defines the behavior of a person that pathfinds between points using the triangulation
type PathfinderBehavior struct {
	Triangulation *Triangulation
	CurrentTarget pixel.Vec
	PathBehavior  *PathBehavior
	Obstacles     []*Obstacle
	TimeWaited    float64
}

// NewPathfinderBehavior creates a new pathfinder behavior.
func NewPathfinderBehavior(triangulation *Triangulation, obstacles []*Obstacle) *PathfinderBehavior {
	return &PathfinderBehavior{
		Triangulation: triangulation,
		CurrentTarget: pixel.Vec{},
		PathBehavior:  NewPathBehavior(nil),
		Obstacles:     obstacles,
		TimeWaited:    0,
	}
}

// GetTarget gets the target of the behavior.
func (b *PathfinderBehavior) GetTarget(p *Person, dt float64) pixel.Vec {
	b.TimeWaited += dt
	if b.CurrentTarget == pixel.ZV || (b.TimeWaited >= 60 && !b.PathBehavior.GoalBehavior.Arrived()) || (b.PathBehavior.GoalBehavior.HasLoitered() && b.PathBehavior.Path.Empty()) {
		b.CurrentTarget = b.Triangulation.Points()[rand.Intn(len(b.Triangulation.Points()))]
		b.PathBehavior.SetPath(AStar(p.Position, b.CurrentTarget, b.Triangulation, b.Obstacles))
		b.TimeWaited = 0
		p.timeSinceLastGoal = 0
	}
	return b.PathBehavior.GetTarget(p, dt)
}

// AStar finds a path between two points using the A* algorithm.
func AStar(start, end pixel.Vec, triangulation *Triangulation, obstacles []*Obstacle) *Path {
	open := []pixel.Vec{}
	cameFrom := map[pixel.Vec]pixel.Vec{}

	var closestToStart pixel.Vec
	for _, v := range triangulation.Points() {
		if start.To(v).Len() < start.To(closestToStart).Len() || closestToStart == pixel.ZV {
			closestToStart = v
		}
	}
	open = append(open, closestToStart)

	gScore := map[pixel.Vec]float64{}
	gScore[closestToStart] = 0

	fScore := map[pixel.Vec]float64{}
	fScore[closestToStart] = end.To(closestToStart).Len()

	for len(open) > 0 {
		current := getLowestFCost(open, fScore)
		if current.To(end).Len() < 10 {
			return reconstructPath(cameFrom, end)
		}
		for i, v := range open {
			if v == current {
				open = append(open[:i], open[i+1:]...)
				break
			}
		}
		for _, v := range triangulation.GetConnectingPoints(current) {
			if lineCollidesObstacles(current, v, obstacles) {
				continue
			}
			t_gScore := gScore[current] + current.To(v).Len()
			g, ok := gScore[v]
			if !ok {
				g = math.Inf(1)
				gScore[v] = g
			}
			if t_gScore < g {
				cameFrom[v] = current
				gScore[v] = t_gScore
				fScore[v] = t_gScore + end.To(v).Len()
				if !slices.Contains(open, v) {
					open = append(open, v)
				}
			}

		}
	}
	panic("No path!")
}

func reconstructPath(cameFrom map[pixel.Vec]pixel.Vec, current pixel.Vec) *Path {
	path := NewPath([]*Goal{NewGoal(current, 100, random(10, 60))})
	next := current

	for {
		v, ok := cameFrom[next]
		if !ok {
			break
		}
		path.goals = append([]*Goal{NewGoal(v, 25, 0)}, path.goals...)
		next = v
	}
	return path
}

func getLowestFCost(open []pixel.Vec, fScore map[pixel.Vec]float64) pixel.Vec {
	var lowest pixel.Vec
	var lowestScore float64 = math.Inf(1)
	for _, v := range open {
		if lowestScore > fScore[v] {
			lowest = v
			lowestScore = fScore[v]
		}
	}
	return lowest
}
