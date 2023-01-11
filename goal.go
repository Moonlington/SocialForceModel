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

// // GoalGraph defines the graph of goals.
// type GoalGraph struct {
// 	Goals       []*Goal
// 	Connections []*GoalEdge
// }

// func NewGoalGraph() *GoalGraph {
// 	return &GoalGraph{}
// }

// func NewGoalGraphFromTriangulation(tri *Triangulation, radius, loiter float64) *GoalGraph {
// 	g := NewGoalGraph()
// 	for _, t := range tri.triangles {
// 		for _, e := range t.Edges() {
// 			goalA := NewGoal(e.A, radius, loiter)
// 			goalB := NewGoal(e.B, radius, loiter)
// 			g.AddConnection(goalA, goalB)
// 			if !g.IncludesGoal(goalA) {
// 				g.AddGoal(goalA)
// 			}
// 			if !g.IncludesGoal(goalB) {
// 				g.AddGoal(goalB)
// 			}
// 		}
// 	}
// 	return g
// }

// // IncludesGoal checks if a goal is included in the graph already.
// func (g *GoalGraph) IncludesGoal(goal *Goal) bool {
// 	for _, g := range g.Goals {
// 		if g == goal {
// 			return true
// 		}
// 	}
// 	return false
// }

// // AddGoal adds a goal to the graph.
// func (g *GoalGraph) AddGoal(goal *Goal) {
// 	g.Goals = append(g.Goals, goal)
// }

// // AddConnection adds a connection between two goals.
// func (g *GoalGraph) AddConnection(a, b *Goal) {
// 	g.Connections = append(g.Connections, NewGoalEdge(a, b))
// }

// // GetConnectingGoals returns the goals that are connected to the given goal.
// func (g *GoalGraph) GetConnectingGoals(goal *Goal) []*Goal {
// 	var goals []*Goal
// 	for _, c := range g.Connections {
// 		if c.A == goal {
// 			goals = append(goals, c.B)
// 		} else if c.B == goal {
// 			goals = append(goals, c.A)
// 		}
// 	}
// 	return goals
// }

// Path describes a list of goals in a specific order
type Path struct {
	goals []*Goal
}

func NewPath(goals []*Goal) *Path {
	return &Path{goals: goals}
}

// GetNextGoal returns the next goal in the path.
func (p *Path) GetNextGoal() *Goal {
	if len(p.goals) == 0 {
		return nil
	}
	returnValue := p.goals[0]
	p.goals = p.goals[1:]
	return returnValue
}

// GetCurrentGoal returns the current goal in the path.
func (p *Path) GetCurrentGoal() *Goal {
	if len(p.goals) == 0 {
		return nil
	}
	return p.goals[0]
}

// GetGoals returns the goals in the path.
func (p *Path) GetGoals() []*Goal {
	return p.goals
}

// Empty returns if the path is empty.
func (p *Path) Empty() bool {
	return len(p.goals) == 0
}

// GetFinishGoal returns the last goal in the path.
func (p *Path) GetFinishGoal() *Goal {
	return p.goals[len(p.goals)-1]
}
