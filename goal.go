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
