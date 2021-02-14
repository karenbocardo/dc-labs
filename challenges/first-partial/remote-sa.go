package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

// Point is a structure that represents a point on a 2D plane
type Point struct {
	X, Y float64
}

// Vector is a structure that represents a vector with magnitude and direction
type Vector struct {
	X, Y float64
}

// Turns two Points into a Vector
func toVector(a, b Point) Vector {
	return Vector{b.X - a.X, b.Y - a.Y}
}

// Calculates cross product between two Vectors
func crossProduct(a, b Vector) float64 {
	return (a.X * b.Y) - (a.Y * b.X)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

//generatePoints array
func generatePoints(s string) ([]Point, error) {
	points := []Point{}

	s = strings.Replace(s, "(", "", -1)
	s = strings.Replace(s, ")", "", -1)
	vals := strings.Split(s, ",")
	if len(vals) < 2 {
		return []Point{}, fmt.Errorf("Point [%v] was not well defined", s)
	}

	var x, y float64

	for idx, val := range vals {

		if idx%2 == 0 {
			x, _ = strconv.ParseFloat(val, 64)
		} else {
			y, _ = strconv.ParseFloat(val, 64)
			points = append(points, Point{x, y})
		}
	}
	return points, nil
}

// Checks if point q lies on line segment pr by the three given collinear points
func onSegment(p, q, r Point) bool {
	if q.X <= math.Max(p.X, r.X) && q.X >= math.Min(p.X, r.X) &&
		q.Y <= math.Max(p.Y, r.Y) && q.Y >= math.Min(p.Y, r.Y) {
		return true
	}
	return false
}

// Finds orientation of a triplet (p on the middle) by using cross product
func orientation(p, q, r Point) int {
	pq := toVector(p, q)
	pr := toVector(p, r)
	cross := crossProduct(pq, pr)
	if cross > 0 {
		return -1 // counterclockwise
	} else if cross < 0 {
		return 1 // clockwise
	} else {
		return 0 // collinear
	}
}

// Finds if two segments formed by four points intersect
func doIntersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)
	if o1 != o2 && o3 != o4 {
		return true
	} else if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	} else if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	} else if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	} else if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}
	return false
}

// Checks wether a list of vertices of a shape have a collision between them
func hasCollision(points []Point) bool {
	points = append(points, points[0])
	nPoints := len(points)
	for i := 0; i < nPoints-1; i++ {
		p1, q1 := points[i], points[i+1]
		for j := i + 2; j < nPoints-1; j++ {
			if i == 0 && j == nPoints-2 {
				continue
			}
			p2, q2 := points[j], points[j+1]
			if doIntersect(p1, q1, p2, q2) {
				return true
			}
		}
	}
	return false
}

// getArea gets the area inside from a given shape
func getArea(points []Point) float64 {
	// shoelace algorithm
	// add first point at the end
	points = append(points, points[0])
	sum1, sum2 := 0.0, 0.0
	for i := 0; i < len(points)-1; i++ {
		sum1 += points[i].X * points[i+1].Y
		sum2 += points[i].Y * points[i+1].X
	}
	return math.Abs(sum1-sum2) / 2.0
}

// Distance between two Points
func getDistance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// getPerimeter gets the perimeter from a given array of connected points
func getPerimeter(points []Point) float64 {
	points = append(points, points[0])
	perimeter := 0.0
	for i := 0; i < len(points)-1; i++ {
		distance := getDistance(points[i], points[i+1])
		perimeter += distance
	}
	return perimeter
}

// handler handles the web request and reponds it
func handler(w http.ResponseWriter, r *http.Request) {

	var vertices []Point
	for k, v := range r.URL.Query() {
		if k == "vertices" {
			points, err := generatePoints(v[0])
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("error: %v", err))
				return
			}
			vertices = points
			break
		}
	}

	// Logging in the server side
	log.Printf("Received vertices array: %v", vertices)

	nVertices := len(vertices)

	response := fmt.Sprintf("Welcome to the Remote Shapes Analyzer\n")
	response += fmt.Sprintf(" - Your figure has : [%v] vertices\n", nVertices)

	if nVertices > 2 {
		collision := false
		if nVertices > 3 {
			// check collisions
			if hasCollision(vertices) {
				collision = true
			}
		}
		if nVertices == 3 || !collision {
			// Results gathering
			area := getArea(vertices)
			perimeter := getPerimeter(vertices)

			// Response construction
			response += fmt.Sprintf(" - Vertices        : %v\n", vertices)
			response += fmt.Sprintf(" - Perimeter       : %v\n", perimeter)
			response += fmt.Sprintf(" - Area            : %v\n", area)
		}
		if collision {
			response += fmt.Sprint("ERROR - Your shape has a collision between some lines.\n")
		}
	} else {
		response += fmt.Sprint("ERROR - Your shape is not compliying with the minimum number of vertices.\n")
	}
	// Send response to client
	fmt.Fprintf(w, response)
}
