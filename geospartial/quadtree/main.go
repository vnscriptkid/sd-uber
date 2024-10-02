package main

import (
	"fmt"
	"math"
)

// Point represents a user location in 2D space.
type Point struct {
	X float64
	Y float64
	// Additional user data can be added here.
}

// Rectangle represents an axis-aligned rectangle in 2D space.
type Rectangle struct {
	X     float64 // Center X-coordinate
	Y     float64 // Center Y-coordinate
	HalfW float64 // Half of the rectangle's width
	HalfH float64 // Half of the rectangle's height
}

// Contains checks if a point is within the rectangle.
func (r *Rectangle) Contains(p *Point) bool {
	// Check if the point's coordinates are within the rectangle's boundaries
	return p.X >= (r.X-r.HalfW) && p.X <= (r.X+r.HalfW) &&
		p.Y >= (r.Y-r.HalfH) && p.Y <= (r.Y+r.HalfH)
}

// Intersects checks if this rectangle intersects with another rectangle.
func (r *Rectangle) Intersects(other *Rectangle) bool {
	// Check if the rectangles overlap by comparing their edges
	return !(other.X-other.HalfW > r.X+r.HalfW ||
		other.X+other.HalfW < r.X-r.HalfW ||
		other.Y-other.HalfH > r.Y+r.HalfH ||
		other.Y+other.HalfH < r.Y-r.HalfH)
}

// Quadtree represents the quadtree node.
type Quadtree struct {
	Boundary *Rectangle
	Capacity int
	Points   []*Point
	Divided  bool
	NW       *Quadtree // Northwest child
	NE       *Quadtree // Northeast child
	SW       *Quadtree // Southwest child
	SE       *Quadtree // Southeast child
}

// NewQuadtree creates a new quadtree node.
func NewQuadtree(boundary *Rectangle, capacity int) *Quadtree {
	return &Quadtree{
		Boundary: boundary,
		Capacity: capacity,
		Points:   make([]*Point, 0),
		Divided:  false,
	}
}

// Insert adds a point to the quadtree.
func (qt *Quadtree) Insert(p *Point) bool {
	// Check if the point is within the quadtree's boundary
	if !qt.Boundary.Contains(p) {
		return false
	}

	// If there's still capacity, add the point to this node
	if len(qt.Points) < qt.Capacity {
		qt.Points = append(qt.Points, p)
		return true
	} else {
		// If the node is at capacity, subdivide (if not already) and insert into children
		if !qt.Divided {
			qt.Subdivide()
		}
		// Try to insert the point into one of the child nodes
		return qt.NW.Insert(p) || qt.NE.Insert(p) || qt.SW.Insert(p) || qt.SE.Insert(p)
	}
}

// Subdivide splits the quadtree node into four children.
func (qt *Quadtree) Subdivide() {
	x := qt.Boundary.X
	y := qt.Boundary.Y
	hw := qt.Boundary.HalfW / 2
	hh := qt.Boundary.HalfH / 2

	// Create four child nodes
	qt.NW = NewQuadtree(&Rectangle{X: x - hw, Y: y - hh, HalfW: hw, HalfH: hh}, qt.Capacity)
	qt.NE = NewQuadtree(&Rectangle{X: x + hw, Y: y - hh, HalfW: hw, HalfH: hh}, qt.Capacity)
	qt.SW = NewQuadtree(&Rectangle{X: x - hw, Y: y + hh, HalfW: hw, HalfH: hh}, qt.Capacity)
	qt.SE = NewQuadtree(&Rectangle{X: x + hw, Y: y + hh, HalfW: hw, HalfH: hh}, qt.Capacity)

	qt.Divided = true
}

// QueryCircle finds all points within a given radius from a center point.
func (qt *Quadtree) QueryCircle(center *Point, radius float64, found *[]*Point) {
	// Create a rectangle that bounds the circle for preliminary intersection testing.
	rangeRect := &Rectangle{
		X:     center.X,
		Y:     center.Y,
		HalfW: radius,
		HalfH: radius,
	}

	// If the query area doesn't intersect this node, return immediately
	if !qt.Boundary.Intersects(rangeRect) {
		return
	}

	// Check all points in this node
	for _, p := range qt.Points {
		// Calculate the distance between the query center and the point
		// [Sqrt](p*p + q*q)
		distance := math.Hypot(p.X-center.X, p.Y-center.Y)
		if distance <= radius {
			*found = append(*found, p)
		}
	}

	// If this node is divided, recursively query the child nodes
	if qt.Divided {
		qt.NW.QueryCircle(center, radius, found)
		qt.NE.QueryCircle(center, radius, found)
		qt.SW.QueryCircle(center, radius, found)
		qt.SE.QueryCircle(center, radius, found)
	}
}

func main() {
	// Define the boundary of the quadtree (e.g., the entire map area).
	boundary := &Rectangle{X: 0, Y: 0, HalfW: 100, HalfH: 100}

	// Create a quadtree with a capacity of 4 points per node.
	qt := NewQuadtree(boundary, 4)

	// Insert some user points into the quadtree.
	points := []*Point{
		{X: -50, Y: -50},
		{X: -40, Y: -40},
		{X: -30, Y: -30},
		{X: -20, Y: -20},
		{X: -10, Y: -10},
		{X: 0, Y: 0},
		{X: 10, Y: 10},
		{X: 20, Y: 20},
		{X: 30, Y: 30},
		{X: 40, Y: 40},
		{X: 50, Y: 50},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Find all users within a radius of 25 units from the point (0, 0).
	center := &Point{X: 0, Y: 0}
	radius := 25.0
	found := make([]*Point, 0)
	qt.QueryCircle(center, radius, &found)

	// Print the results
	fmt.Printf("Found %d users within radius %.2f of point (%.2f, %.2f):\n", len(found), radius, center.X, center.Y)
	for _, p := range found {
		fmt.Printf("User at (%.2f, %.2f)\n", p.X, p.Y)
	}
}
