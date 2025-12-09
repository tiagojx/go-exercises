package main

import "fmt"

type User struct {
	ID        int
	FirstName string
	LastName  string
	IsActive  bool
}

// the function inherdit the struct like a object-oriented paradigm.
// a pointer is not needed because we'll return a new 'ghost' variable.
func (u User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// also works with methods.
// a pointer is needed because we'll change the variable value directely.
func (u *User) Deactivate() {
	u.IsActive = false
}

// also works in toggle mode
func (u *User) Toggle() {
	u.IsActive = !u.IsActive
}

func main() {
	// instantiation of a struct in Go.
	u := User{ID: 1, FirstName: "Cleitin", LastName: "Dischava", IsActive: true}

	fmt.Println(u.FullName())
	fmt.Println("Ativo: ", u.IsActive)

	// deactivate
	fmt.Println("Desativando usuário...")
	u.Deactivate()
	fmt.Println("Ativo: ", u.IsActive)

	// toggle
	fmt.Println("Switch - ativar/desativar:")
	u.Toggle()
	fmt.Println("Ativo: ", u.IsActive)
	u.Toggle()
	fmt.Println("Ativo: ", u.IsActive)
}
