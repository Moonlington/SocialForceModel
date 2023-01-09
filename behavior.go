package main

import (
	"math/rand"

	"github.com/faiface/pixel"
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

// ChooseNextWanderLocation chooses the next wander location.
func (b *WanderBehavior) ChooseNextWanderLocation(p *Person) *Goal {
	var possibleGoals []*Goal
	for _, goal := range b.WanderGoals {
		intersects := false
		for _, obstacle := range b.Obstacles {
			if obstacle.Inner {
				continue
			}
			if obstacle.IntersectLine(pixel.L(p.Position, goal.Target)).Len() > 0 {
				intersects = true
				break
			}
		}
		if !intersects {
			possibleGoals = append(possibleGoals, goal)
		}
	}
	if len(possibleGoals) == 0 {
		return nil
	}
	return possibleGoals[rand.Intn(len(possibleGoals))]
}

type PathfinderBehavior struct {
	goalGraph   *GoalGraph
	currentPath *Path
	obstacles   []*Obstacle
}

func NewPathfinderBehavior(goalGraph *GoalGraph) *PathfinderBehavior {
	return &PathfinderBehavior{goalGraph: goalGraph}
}

// SelectRandomGoal randomly selects a goal from the goal graph.
func (b *PathfinderBehavior) SelectRandomGoal() *Goal {
	return b.goalGraph.Goals[rand.Intn(len(b.goalGraph.Goals))]
}

// PathfindToGoal uses A* algorithm to find a path to a given goal from where the person is located
func (b *PathfinderBehavior) PathfindToGoal(p *Person, goal *Goal) *Path {
	if b.currentPath != nil && b.currentPath.GetFinishGoal() == goal {
		return b.currentPath
	}
	return b.currentPath
}
