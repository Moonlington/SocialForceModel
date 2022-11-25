package main

import (
	"image/color"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Person struct {
	id    int
	Color color.RGBA

	Position pixel.Vec
	Velocity pixel.Vec

	Goals       []*Goal
	CurrentGoal int
	Loitered    float64

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
	p.Color = colornames.Cyan

	p.Position = pixel.V(0, 0)
	p.Velocity = pixel.V(0, 0)

	p.Goals = []*Goal{newGoal(pixel.V(500, 350), 100, 0)}
	p.CurrentGoal = 0
	p.Loitered = 0.

	p.DesiredSpeed = 13.3
	p.Mass = 60.
	p.alpha = 1.

	p.wallThreshold = 1.

	p.Radius = 0.2

	return p
}

func (p *Person) willForce(dt float64) pixel.Vec {
	gw := p.Mass * p.alpha

	if len(p.Goals) <= p.CurrentGoal {
		return pixel.V(0, 0).Sub(p.Velocity).Scaled(gw)
	}

	target := p.Goals[p.CurrentGoal].Target
	Vd := p.Position.To(target).Unit().Scaled(p.DesiredSpeed)

	if p.Position.To(target).Len() < p.Goals[p.CurrentGoal].Range {
		Vd = pixel.V(0, 0)
		if p.Goals[p.CurrentGoal].LoiterAfter >= p.Loitered {
			p.Loitered += dt
		} else {
			p.Loitered = 0
			p.CurrentGoal++
		}
	}
	return Vd.Sub(p.Velocity).Scaled(gw)
}

func (p *Person) intermediateRangeForce(o *Person) pixel.Vec {
	fmax := p.Mass * 2. * p.alpha

	t := p.Velocity.Unit()
	n := p.Velocity.Normal().Unit()

	dTm := -(p.Position.Sub(o.Position).Dot(p.Velocity.Sub(o.Velocity)) / p.Velocity.Sub(o.Velocity).Dot(p.Velocity.Sub(o.Velocity)))
	if dTm < 0 || dTm > 10 {
		return pixel.V(0, 0)
	}

	// rhot := t.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)
	// rhon := n.Scaled(p.Position.Sub(o.Position).Len()).Len() / (p.Radius)

	rhot := p.Position.Sub(o.Position).Project(t).Len() / (p.Radius)
	rhon := p.Position.Sub(o.Position).Project(n).Len() / (p.Radius)

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
	minDistVec := pixel.V(math.Inf(1), math.Inf(1))

	for _, o := range obstacles {
		d := o.Dist(p)
		if minDistVec.Len() > d.Len() {
			minDistVec = d
		}
	}

	if minDistVec.Len() > p.wallThreshold {
		return pixel.V(0, 0)
	}

	fmax := p.Mass * 16. * p.alpha
	s := minDistVec.Unit()
	return s.Scaled(-fmax * (1 / (1 + math.Pow(minDistVec.Len()/p.Radius, 2))))
}

func (p *Person) motionInhibition(obstacles []*Obstacle) {
	minDistVec := pixel.V(math.Inf(1), math.Inf(1))

	for _, o := range obstacles {
		d := o.Dist(p)
		if minDistVec.Len() > d.Len() {
			minDistVec = d
		}
	}

	if minDistVec.Len() > p.Radius || p.Velocity.Dot(minDistVec) <= 1 {
		return
	}

	// p.sumForce = p.sumForce.Sub(p.sumForce.Project(minDistVec))
	p.Velocity = p.Velocity.Project(minDistVec.Normal())
}

func (p *Person) fixCollision(obstacles []*Obstacle) {
	minDistVec := pixel.V(math.Inf(1), math.Inf(1))
	var closestObstacle Obstacle

	for _, o := range obstacles {
		d := o.Dist(p)
		if minDistVec.Len() > d.Len() {
			minDistVec = d
			closestObstacle = *o
		}
	}

	if minDistVec.Len() > p.Radius*.9 {
		return
	}

	p.Position = p.Position.Add(pixel.C(p.Position, p.Radius).IntersectRect(closestObstacle.Rect))
}

func (p *Person) update(dt float64, others []*Person, obstacles []*Obstacle) {
	p.sumForce = pixel.V(0, 0)

	p.sumForce = p.sumForce.Add(p.willForce(dt))
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		p.sumForce = p.sumForce.Add(p.intermediateRangeForce(o))
		p.sumForce = p.sumForce.Add(p.nearRangeForce(o))
		p.sumForce = p.sumForce.Add(p.contactForce(o))
	}
	p.sumForce = p.sumForce.Add(p.wallForce(obstacles))

	// if p.id == 0 {
	// 	fmt.Println(p.sumForce)
	// }

	p.fixCollision(obstacles)
	p.Velocity = p.Velocity.Add(p.sumForce.Scaled(1 / p.Mass).Scaled(dt))
	p.motionInhibition(obstacles)
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))

}

func (p *Person) Draw(imd *imdraw.IMDraw) {
	imd.Color = p.Color
	imd.Push(p.Position)
	imd.Circle(p.Radius, 1)

	imd.Color = colornames.Lime
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.Velocity))
	imd.Line(1)

	imd.Color = colornames.Yellow
	imd.Push(p.Position)
	imd.Push(p.Position.Add(p.sumForce.Scaled(1 / p.Mass)))
	imd.Line(1)

	// imd.Color = colornames.Magenta
	// imd.Push(p.Position)
	// imd.Circle(p.wallThreshold, 1)

	// if len(p.Goals) > p.CurrentGoal {
	// 	imd.Color = colornames.Red
	// 	imd.Push(p.Position)
	// 	imd.Push(p.Goals[p.CurrentGoal].Target)
	// 	imd.Line(1)
	// }
}
