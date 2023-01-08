package main

import "github.com/faiface/pixel"

// Goal defines the goal of a person.
type Goal struct {
	Target      pixel.Vec
	Range       float64
	LoiterAfter float64
}

func NewGoal(target pixel.Vec, r, loiter float64) *Goal {
	return &Goal{Target: target, Range: r, LoiterAfter: loiter}
}

// GoalEdge defines the edge between two goals.
type GoalEdge struct {
	A *Goal
	B *Goal
}

func NewGoalEdge(a, b *Goal) *GoalEdge {
	return &GoalEdge{A: a, B: b}
}

// GoalGraph defines the graph of goals.
type GoalGraph struct {
	Goals       []*Goal
	Connections []*GoalEdge
}

func NewGoalGraph() *GoalGraph {
	return &GoalGraph{}
}

func NewGoalGraphFromTriangulation(tri *Triangulation, radius, loiter float64) *GoalGraph {
	g := NewGoalGraph()
	for _, t := range tri.triangles {
		for _, p := range t.points {
			g.AddGoal(NewGoal(p, radius, loiter))
		}
	}
	for _, t := range tri.triangles {
		for _, e := range t.Edges() {

		}
	}
	return g
}

// AddGoal adds a goal to the graph.
func (g *GoalGraph) AddGoal(goal *Goal) {
	g.Goals = append(g.Goals, goal)
}

// AddConnection adds a connection between two goals.
func (g *GoalGraph) AddConnection(a, b *Goal) {
	g.Connections = append(g.Connections, NewGoalEdge(a, b))
}

// GetConnectingGoals returns the goals that are connected to the given goal.
func (g *GoalGraph) GetConnectingGoals(goal *Goal) []*Goal {
	var goals []*Goal
	for _, c := range g.Connections {
		if c.A == goal {
			goals = append(goals, c.B)
		} else if c.B == goal {
			goals = append(goals, c.A)
		}
	}
	return goals
}
