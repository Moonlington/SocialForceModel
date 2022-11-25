package main

import "github.com/faiface/pixel"

type Goal struct {
	Target pixel.Vec

	Range float64

	LoiterAfter float64
}

func newGoal(target pixel.Vec, r, loiter float64) *Goal {
	g := new(Goal)

	g.Target = target
	g.Range = r
	g.LoiterAfter = loiter

	return g
}
