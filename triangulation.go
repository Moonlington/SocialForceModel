package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"golang.org/x/image/colornames"
)

type Triangle struct {
	points [3]pixel.Vec
}

func NewTriangle(p1, p2, p3 pixel.Vec) *Triangle {
	return &Triangle{
		points: [3]pixel.Vec{p1, p2, p3},
	}
}

func (T *Triangle) Edges() [3]pixel.Line {
	return [3]pixel.Line{
		pixel.L(T.points[0], T.points[1]),
		pixel.L(T.points[1], T.points[2]),
		pixel.L(T.points[2], T.points[0]),
	}
}

func (T *Triangle) Draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.Darkorange
	imd.Push(T.points[0], T.points[1], T.points[2])
	imd.Polygon(1)
	imd.Color = colornames.White
	imd.Push(T.points[0])
	imd.Circle(2, 0)
	imd.Push(T.points[1])
	imd.Circle(2, 0)
	imd.Push(T.points[2])
	imd.Circle(2, 0)
	// imd.Color = colornames.Darkorange
	// imd.Push(T.Circumcircle().Center)
	// imd.Circle(T.Circumcircle().Radius, 2)
}

func (T *Triangle) Circumcircle() pixel.Circle {
	a := T.points[0]
	b := T.points[1]
	c := T.points[2]

	A := a.X*a.X + a.Y*a.Y
	B := b.X*b.X + b.Y*b.Y
	C := c.X*c.X + c.Y*c.Y

	D := a.X*(b.Y-c.Y) + b.X*(c.Y-a.Y) + c.X*(a.Y-b.Y)

	x := (A*(b.Y-c.Y) + B*(c.Y-a.Y) + C*(a.Y-b.Y)) / (2 * D)
	y := (A*(c.X-b.X) + B*(a.X-c.X) + C*(b.X-a.X)) / (2 * D)

	r := pixel.V(x, y).Sub(a).Len()

	return pixel.C(pixel.V(x, y), r)
}

type Triangulation struct {
	triangles []*Triangle
}

func BowyerWatson(points []pixel.Vec) *Triangulation {
	T := new(Triangulation)
	// Add super triangle
	T.AddTriangle(NewTriangle(pixel.V(-10000, -10000), pixel.V(10000, -10000), pixel.V(0, 10000)))
	// Add points
	for _, p := range points {
		// fmt.Println(p)
		badTriangles := make([]*Triangle, 0)
		for _, t := range T.triangles {
			if t.Circumcircle().Contains(p) {
				badTriangles = append(badTriangles, t)
			}
		}
		T.RemoveTriangles(badTriangles)
		// fmt.Println(badTriangles)
		polygon := make([]pixel.Vec, 0)
		for _, t := range badTriangles {
			for _, e := range t.Edges() {
				shared := false
				for _, t2 := range badTriangles {
					if t == t2 {
						continue
					}
					for _, e2 := range t2.Edges() {
						if e == e2 || (e.A == e2.B && e.B == e2.A) {
							shared = true
							break
						}
					}
					if shared {
						break
					}
				}
				if !shared {
					polygon = append(polygon, e.A, e.B)
				}
			}
		}
		// fmt.Println(polygon)
		for i := 0; i < len(polygon); i += 2 {
			T.AddTriangle(NewTriangle(polygon[i], polygon[i+1], p))
		}
	}
	// Remove super triangle
	remove := make([]*Triangle, 0)
	for _, t := range T.triangles {
		for _, p := range t.points {
			if p.X == -10000 || p.X == 10000 || p.Y == 10000 {
				remove = append(remove, t)
				break
			}
		}
	}
	T.RemoveTriangles(remove)
	return T
}

func (T *Triangulation) AddTriangle(t *Triangle) {
	T.triangles = append(T.triangles, t)
}

func (T *Triangulation) AddTriangles(triangles []*Triangle) {
	T.triangles = append(T.triangles, triangles...)
}

func (T *Triangulation) RemoveTriangle(t *Triangle) {
	for i, t2 := range T.triangles {
		if t == t2 {
			T.triangles = append(T.triangles[:i], T.triangles[i+1:]...)
			return
		}
	}
	panic("Triangle not found")
}

func (T *Triangulation) RemoveTriangles(triangles []*Triangle) {
	for _, t := range triangles {
		T.RemoveTriangle(t)
	}
}

func (T *Triangulation) Draw(imd *imdraw.IMDraw) {
	for _, t := range T.triangles {
		t.Draw(imd)
	}
}

func (T *Triangulation) Points() []pixel.Vec {
	var points []pixel.Vec
	for _, t := range T.triangles {
		for _, p := range t.points {
			found := false
			for _, p2 := range points {
				if p == p2 {
					found = true
					break
				}
			}
			if !found {
				points = append(points, p)
			}
		}
	}
	return points
}

func (T *Triangulation) Edges() []pixel.Line {
	var edges []pixel.Line
	for _, t := range T.triangles {
		for _, e := range t.Edges() {
			found := false
			for _, e2 := range edges {
				if e == e2 || (e.A == e2.B && e.B == e2.A) {
					found = true
					break
				}
			}
			if !found {
				edges = append(edges, e)
			}
		}
	}
	return edges
}

func (T *Triangulation) GetConnectingPoints(p pixel.Vec) []pixel.Vec {
	var points []pixel.Vec
	for _, e := range T.Edges() {
		if e.A == p {
			points = append(points, e.B)
		} else if e.B == p {
			points = append(points, e.A)
		}
	}
	return points
}
