package main

import (
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

	wallThreshold float64
}

func newPerson(id int) *Person {
	p := new(Person)

	p.id = id

	p.Position = pixel.V(0, 0)
	p.Velocity = pixel.V(0, 0)
	p.Goal = pixel.V(500, 350)
	p.DesiredSpeed = 13.3
	p.Mass = 60.
	p.alpha = 1.

	p.wallThreshold = 1.

	p.Radius = 0.2

	return p
}

func (p *Person) willForce() pixel.Vec {
	gw := p.Mass * p.alpha
	Vd := p.Position.To(p.Goal).Unit().Scaled(p.DesiredSpeed)
	if p.Position.To(p.Goal).Len() <= 2*p.Radius {
		Vd = Vd.Scaled(1 / (1 + (p.Radius / p.Position.To(p.Goal).Len())))
	}
	return Vd.Sub(p.Velocity).Scaled(gw)
}

func (p *Person) intermediateRangeForce(o *Person) pixel.Vec {
	fmax := p.Mass * 2. * p.alpha

	t := p.Velocity.Unit()
	n := p.Velocity.Normal().Unit()

	rhot := t.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)
	rhon := n.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)

	return t.Scaled(-fmax * (1 / (1 + math.Pow(rhot, 2)))).Add(n.Scaled(-fmax * (1 / (1 + math.Pow(rhon, 2)))))
}

func (p *Person) nearRangeForce(o *Person) pixel.Vec {
	fmax := p.Mass * 4. * p.alpha
	rho := p.Position.Sub(o.Position).Len() / (p.Radius)
	return p.Position.To(o.Position).Unit().Scaled(-fmax * (1 / (1 + math.Pow(rho, 2))))
}

func (p *Person) contactForce(o *Person) pixel.Vec {
	sumForce := pixel.V(0, 0)
	rho := p.Position.Sub(o.Position).Len() / (p.Radius + o.Radius)

	fmax := p.Mass * 8. * math.Max(p.alpha, o.alpha)
	var f pixel.Vec
	if rho <= 1 {
		f = p.Position.To(o.Position).Unit().Scaled(-2 * fmax * (1 / (1 + math.Pow(rho, 2))))
	} else {
		f = p.Position.To(o.Position).Unit().Scaled(-1 * fmax * (1 / (1 + math.Pow(rho, 2))))
	}

	sumForce = sumForce.Add(f)
	t := p.Position.To(o.Position).Unit().Normal()
	var ft pixel.Vec
	ft = t.Scaled(0.2 * f.Len() * 1)
	sumForce = sumForce.Add(ft)
	return sumForce
}

func (p *Person) wallForce(obstacles []*Obstacle) pixel.Vec {
	var minDistObstacle int
	minDist := math.Inf(1)

	for i, o := range obstacles {
		d := o.Dist(p).Len()
		if minDist > d {
			minDist = d
			minDistObstacle = i
		}
	}

	fmax := p.Mass * 16. * p.alpha

	s := obstacles[minDistObstacle].Dist(p).Unit()

	return s.Scaled(-fmax * (1 / (1 + math.Pow(minDist/p.Radius, 2))))
}

func (p *Person) update(dt float64, others []*Person, obstacles []*Obstacle) {
	p.sumForce = pixel.V(0, 0)

	p.sumForce = p.sumForce.Add(p.willForce())
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		p.sumForce = p.sumForce.Add(p.intermediateRangeForce(o))
		p.sumForce = p.sumForce.Add(p.nearRangeForce(o))
		p.sumForce = p.sumForce.Add(p.contactForce(o))
		// p.sumForce = p.sumForce.Add(p.wallForce(obstacles))
	}

	// if p.id == 0 {
	// 	fmt.Println(p.sumForce)
	// }

	p.Velocity = p.Velocity.Add(p.sumForce.Scaled(1 / p.Mass).Scaled(dt))
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))
}

func (p *Person) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.Cyan
	imd.Push(p.Position)
	imd.Circle(p.Radius, 1)

	// imd.Color = colornames.Magenta
	// imd.Push(p.Position)
	// imd.Circle(p.wallThreshold, 1)

	imd.Color = colornames.Lime
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.Velocity))
	imd.Line(1)

	imd.Color = colornames.Yellow
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.sumForce.Scaled(1 / p.Mass)))
	imd.Line(1)

	// imd.Color = colornames.Red
	// imd.Push(p.Position)
	// imd.Push(p.Goal)
	// imd.Line(1)
}
