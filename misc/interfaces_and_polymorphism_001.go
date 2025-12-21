package main

import "fmt"

// define the interface.
// all structs who contains the functions listed here will be
// handled by this intercade.
type Noisy interface {
	makeNoise() string
}

type Dog struct {
	Name string
}

// Noisy's makeNoise() implementation A
func (d Dog) makeNoise() string {
	return d.Name + " diz: au-au!"
}

type Cat struct {
	Name string
}

// Noisy's makeNoise() implementation B
func (c Cat) makeNoise() string {
	return c.Name + " diz: miau!"
}

type Vaca struct {
	Name string
	Raca string
}

func (v Vaca) makeNoise() string {
	return v.Name + " da raça " + v.Raca + " diz: muuuuuuu!"
}

// this method get any object of the Noisy interface.
func petSounds(n Noisy) {
	fmt.Println(n.makeNoise())
}

func main() {
	d := Dog{Name: "Bifin"}
	c := Cat{Name: "Vladmir"}
	v := Vaca{Name: "Mimosa", Raca: "Nelore"}

	petSounds(d)
	petSounds(c)
	petSounds(v)
}
