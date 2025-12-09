package main

import (
	"fmt"
	"math"
)

// define the interface.
// all structs who contains the functions listed here will be
// handled by this intercade.
type Shape interface {
	Area() float64
	GetName() string
}

type Rectangle struct {
	Width  float64
	Height float64
}

func (r Rectangle) Area() float64 {
	return (r.Height * 2) + (r.Width * 2)
}

func (r Rectangle) GetName() string {
	return "Retângulo"
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * (c.Radius * c.Radius)
}

func (c Circle) GetName() string {
	return "Círculo"
}

func main() {
	r := Rectangle{Width: 30, Height: 10}
	c := Circle{Radius: 50}

	shapes := []Shape{}
	shapes = append(shapes, r, c)
	for s := 0; s < len(shapes); s++ {
		fmt.Println(shapes[s].GetName(), ": ", shapes[s].Area())
	}
}
