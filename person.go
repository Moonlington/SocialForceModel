package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Person struct {
	Position     pixel.Vec
	Velocity     pixel.Vec
	Acceleration pixel.Vec
	Goal         pixel.Vec
}

func newPerson() *Person {
	p := new(Person)

	p.Position = pixel.V(0, 0)
	p.Velocity = pixel.V(0, 0)
	p.Acceleration = pixel.V(0, 0)
	p.Goal = pixel.V(500, 300)

	return p
}

func (p *Person) update(dt float64) {
	p.Velocity = p.Velocity.Add(p.Acceleration.Scaled(dt))
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))
	// if p.Position.X > 1024 {
	// 	p.Position.X -= 1024
	// }
	// if p.Position.Y > 768 {
	// 	p.Position.Y -= 768
	// }
	// if p.Position.X < 0 {
	// 	p.Position.X += 1024
	// }
	// if p.Position.Y < 0 {
	// 	p.Position.Y += 768
	// }
}

func (p *Person) draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.Orange
	imd.Push(p.Position)
	imd.Circle(10, 1)
}
