package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const SCALING float64 = 100.

type Person struct {
	id    int
	Color color.RGBA

	Position pixel.Vec
	Velocity pixel.Vec

	Behavior Behavior

	Radius       float64
	DesiredSpeed float64
	Mass         float64

	alpha    float64
	sumForce pixel.Vec

	wallThreshold float64

	timeSinceLastGoal float64
}

func newPerson(id int) *Person {
	p := new(Person)

	p.id = id
	p.Color = colornames.Cyan

	p.Position = pixel.V(0, 0)
	p.Velocity = pixel.V(0, 0)

	p.Behavior = nil

	p.DesiredSpeed = math.Max(0.01, (rand.NormFloat64()*0.025+1.)*SCALING)
	p.Mass = rand.NormFloat64()*5 + 70
	// p.alpha = 1. * math.Sqrt(SCALING)
	p.alpha = 2.

	p.Radius = (rand.NormFloat64()*0.025 + 0.2) * SCALING
	p.wallThreshold = math.Max(p.Radius, (rand.NormFloat64()*.5+1.)*SCALING)

	p.timeSinceLastGoal = 0.

	return p
}

func (p *Person) XY() (float64, float64) {
	return p.Position.XY()
}

func (p *Person) willForce(dt float64, target pixel.Vec) pixel.Vec {
	gw := p.Mass * p.alpha * (1 + p.timeSinceLastGoal/20)
	if target == p.Position {
		p.timeSinceLastGoal = 0
		return pixel.V(0, 0).Sub(p.Velocity).Scaled(gw)
	}
	Vd := p.Position.To(target).Unit().Scaled(p.DesiredSpeed * (1 + p.timeSinceLastGoal/60))
	p.timeSinceLastGoal += dt
	return Vd.Sub(p.Velocity).Scaled(gw)
}

func (p *Person) intermediateRangeForce(o *Person) pixel.Vec {
	fmax := p.Mass * 16. * p.alpha

	t := p.Velocity.Unit()
	n := p.Velocity.Normal().Unit()

	dTm := -(p.Position.Sub(o.Position).Dot(p.Velocity.Sub(o.Velocity)) / p.Velocity.Sub(o.Velocity).Dot(p.Velocity.Sub(o.Velocity)))
	if dTm <= 0 {
		return pixel.V(0, 0)
	}

	rhot := p.Position.Sub(o.Position).Project(t).Len() / (p.Radius)
	rhon := p.Position.Sub(o.Position).Project(n).Len() / (p.Radius)

	return t.Scaled(-fmax * (1 / (1 + math.Pow(rhot, 2)))).Add(n.Scaled(-fmax * (1 / (1 + math.Pow(rhon, 2)))))
}

func (p *Person) nearRangeForce(o *Person) pixel.Vec {
	fmax := p.Mass * 64. * p.alpha
	rho := p.Position.Sub(o.Position).Len() / (p.Radius)
	return p.Position.To(o.Position).Unit().Scaled(-fmax * (1 / (1 + math.Pow(rho, 2))))
}

func (p *Person) contactForce(o *Person) pixel.Vec {
	sumForce := pixel.V(0, 0)
	rho := p.Position.Sub(o.Position).Len() / (p.Radius + o.Radius)

	fmax := p.Mass * 128. * math.Max(p.alpha, o.alpha)
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

	fmax := p.Mass * 256. * p.alpha
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
		if o.Inner {
			continue
		}
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

func (p *Person) fixCollisionOthers(others []*Person) {
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		distance := p.Position.To(o.Position).Len()
		overlap := -(distance - p.Radius*.9 - o.Radius)
		if overlap <= 0 {
			continue
		}
		p.Position = p.Position.Add(p.Position.To(o.Position).Unit().Scaled(-overlap))
	}
}

func (p *Person) kinematicConstraint(dt float64, others []*Person) {
	xi := .75
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		var ds float64
		if p.Velocity.Dot(o.Velocity) > 0 {
			ds = xi*p.Radius + xi*o.Radius
		} else {
			ds = xi*p.Radius + (1-xi)*o.Radius
		}

		pNewPosition := p.Position.Add(p.Velocity.Scaled(dt))
		oNewPosition := o.Position.Add(o.Velocity.Scaled(dt))

		dn := pNewPosition.To(oNewPosition).Len()

		cr := math.Min(1, math.Max(0, (dn-ds)/ds))

		distance := p.Position.To(o.Position)
		p.Velocity = p.Velocity.Project(distance).Scaled(cr).Add(p.Velocity.Project(distance.Normal()))
	}
}

func (p *Person) update(dt float64, others []*Person, obstacles []*Obstacle) {
	p.sumForce = pixel.V(0, 0)

	p.sumForce = p.sumForce.Add(p.willForce(dt, p.Behavior.GetTarget(p, dt)))
	for _, o := range others {
		if o.id == p.id {
			continue
		}
		p.sumForce = p.sumForce.Add(p.intermediateRangeForce(o))
		p.sumForce = p.sumForce.Add(p.nearRangeForce(o))
		p.sumForce = p.sumForce.Add(p.contactForce(o))
	}
	p.sumForce = p.sumForce.Add(p.wallForce(obstacles))

	p.fixCollisionOthers(others)
	p.fixCollision(obstacles)
	p.Velocity = p.Velocity.Add(p.sumForce.Scaled(1 / p.Mass).Scaled(dt))
	p.motionInhibition(obstacles)
	p.kinematicConstraint(dt, others[:])
	p.Position = p.Position.Add(p.Velocity.Scaled(dt))

}

func (p *Person) Draw(win *pixelgl.Window, imd *imdraw.IMDraw) {
	if win.MousePosition().Sub(win.Bounds().Center()).To(p.Position).Len() < p.Radius {
		imd.Color = colornames.Red
		imd.Push(p.Position)
		imd.Circle(p.Radius, 1)
		p.DrawGoal(imd)
	} else {
		imd.Color = p.Color
		imd.Push(p.Position)
		imd.Circle(p.Radius, 1)
	}

	// imd.Color = colornames.Lime
	// imd.Push(p.Position)
	// imd.Push(p.Position.Add(p.Velocity))
	// imd.Line(1)

	// imd.Color = colornames.Yellow
	// imd.Push(p.Position)
	// imd.Push(p.Position.Add(p.sumForce.Scaled(1 / p.Mass)))
	// imd.Line(1)

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

// DrawGoal draws a line between the person and the goal
func (p *Person) DrawGoal(imd *imdraw.IMDraw) {
	target := p.Behavior.GetTarget(p, 0)
	imd.Color = colornames.Lime
	imd.Push(p.Position)
	imd.Push(target)
	imd.Line(1)
}
