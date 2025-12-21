package main

import (
	"errors"
	"fmt"
)

// a fundamental divide function; it may return a value (float64) and a error when
// the user trying to divide a by b=0.
func Divide(a, b float64) (float64, error) {
	if b != 0 {
		return (a / b), nil // return de value of the division, and 'nil' for the error
	} else {
		return 0.0, errors.New("não é possível dividir por zero.")
	}
}

func main() {
	var a (float64) = 8
	var b (float64) = 4
	var c (float64) = 0

	// possible to divide
	fmt.Printf("Dividindo %.1f por %.1f... ", a, b)
	value, err := Divide(a, b)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("O resultado da divisão é:", value)
	}

	// impossible to divide
	fmt.Printf("Dividindo %.1f por %.1f... ", a, c)
	value, err = Divide(a, c)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("O resultado da divisão é:", value)
	}
}
