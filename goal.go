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

// Path defines a path of goals.
type Path struct {
	Goals []*Goal
}

func newPath(goals ...*Goal) *Path {
	return &Path{Goals: goals}
}

func (p *Path) AddGoal(goal *Goal) {
	p.Goals = append(p.Goals, goal)
}

func (p *Path) RemoveGoal(goal *Goal) {
	for i, g := range p.Goals {
		if g == goal {
			p.Goals = append(p.Goals[:i], p.Goals[i+1:]...)
			return
		}
	}
}

func (p *Path) NextGoal() *Goal {
	next := p.Goals[0]
	p.Goals = p.Goals[1:]
	return next
}

func (p *Path) IsEmpty() bool {
	return len(p.Goals) == 0
}
