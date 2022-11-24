package main

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Person struct {
	id int

	Position pixel.Vec
	Velocity pixel.Vec

	Goal pixel.Vec

	Radius       float64
	DesiredSpeed float64
	Mass         float64

	alpha    float64
	sumForce pixel.Vec
}

func newPerson(id int) *Person {
	p := new(Person)

	p.id = id

	p.Position = pixel.V(0, 0)
	p.Velocity = pixel.V(0, 0)
	p.Goal = pixel.V(500, 350)
	p.DesiredSpeed = 13.3
	p.Mass = 60.
	p.alpha = 2.

	p.Radius = 10

	return p
}

func (p *Person) checkBoundaries(b pixel.Vec) {
	if p.Position.X-p.Radius <= 0 || p.Position.X+p.Radius >= b.X {
		p.Velocity.X = -p.Velocity.X
	}
	if p.Position.Y-p.Radius <= 0 || p.Position.Y+p.Radius >= b.Y {
		p.Velocity.Y = -p.Velocity.Y
	}
}

func (p *Person) willForce() pixel.Vec {
	gw := math.Min(p.Mass*p.alpha*p.Position.To(p.Goal).Len(), p.Mass*p.alpha)
	Vd := p.Position.To(p.Goal).Unit().Scaled(p.DesiredSpeed)
	f := Vd.Sub(p.Velocity).Scaled(gw)
	return f
}

func (p *Person) intermediateRangeForce(others []*Person) pixel.Vec {
	sumForce := pixel.V(0, 0)
	fmax := p.Mass * 2. * p.alpha
	v1 := p.Velocity
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		t := v1.Unit()
		n := v1.Normal().Unit()

		rhot := t.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)
		rhon := n.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)

		f := t.Scaled(-fmax * (1 / (1 + math.Pow(rhot, 2)))).Add(n.Scaled(-fmax * (1 / (1 + math.Pow(rhon, 2)))))

		sumForce = sumForce.Add(f)
	}
	return sumForce
}

func (p *Person) nearRangeForce(others []*Person) pixel.Vec {
	sumForce := pixel.V(0, 0)
	fmax := p.Mass * 4. * p.alpha
	for _, o := range others {
		if o.id == p.id {
			continue
		}

		rho := p.Position.Sub(o.Position).Len() / (p.Radius)

		f := p.Position.To(o.Position).Unit().Scaled(-fmax * (1 / (1 + math.Pow(rho, 2))))

		sumForce = sumForce.Add(f)
	}
	return sumForce
}

func (p *Person) contactForce(others []*Person) pixel.Vec {
	sumForce := pixel.V(0, 0)
	for _, o := range others {
		if o.id == p.id {
			continue
		}

		rho := p.Position.Sub(o.Position).Len() / (p.Radius + o.Radius)

		fmax := p.Mass * 8. * math.Max(p.alpha, o.alpha)
		var f pixel.Vec
		if rho < 1 {
			f = p.Position.To(o.Position).Unit().Scaled(-fmax * (1 / (1 + math.Pow(rho, 2))))
		} else {
			f = p.Position.To(o.Position).Unit().Scaled(-2 * fmax * (1 / (1 + math.Pow(rho, 2))))
		}

		sumForce = sumForce.Add(f)
		t := p.Position.To(o.Position).Unit().Normal()
		var ft pixel.Vec
		if math.Signbit(t.Dot(p.Velocity.Sub(o.Velocity))) {
			ft = t.Scaled(0.2 * f.Len() * 1)
		} else {
			ft = t.Scaled(0.2 * f.Len() * -1)
		}
		sumForce = sumForce.Add(ft)
	}
	return sumForce
}

func (p *Person) update(dt float64, others []*Person) {
	p.sumForce = pixel.V(0, 0)

	p.sumForce = p.sumForce.Add(p.willForce())
	p.sumForce = p.sumForce.Add(p.intermediateRangeForce(others))
	p.sumForce = p.sumForce.Add(p.nearRangeForce(others))
	p.sumForce = p.sumForce.Add(p.contactForce(others))

	if p.id == 0 {
		fmt.Println(p.sumForce)
	}

	p.Velocity = p.Velocity.Add(p.sumForce.Scaled(1 / p.Mass).Scaled(dt))
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))
}

func (p *Person) draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.Orange
	imd.Push(p.Position)
	// if p.id == 1 {
	imd.Circle(p.Radius, 1)
	imd.Color = colornames.Lime
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.Velocity))
	imd.Line(2)
	imd.Color = colornames.Yellow
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.sumForce.Scaled(1 / p.Mass)))
	imd.Line(2)
	// } else {
	// 	imd.Circle(p.Radius, 1)
	// }
}
