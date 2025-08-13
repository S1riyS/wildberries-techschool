package main

import "fmt"

// Parent structure - Human
type Human struct {
	Name string
	Age  int
}

// Human' method
func (h *Human) SayHello() {
	fmt.Printf("Hello, my name is %s and I'm %d years old.\n", h.Name, h.Age)
}

// Child structure - Actions. Human is embedded
type Action struct {
	Human
	Profession string
}

// Action's method
func (a *Action) Work() {
	fmt.Printf("%s works as %s.\n", a.Name, a.Profession)
}

func main() {
	person := Action{
		Human: Human{
			Name: "Kirill",
			Age:  20,
		},
		Profession: "Go developer",
	}

	// Call Human's method
	person.SayHello()

	// Call Action's method
	person.Work()

	// Access fields
	fmt.Println("Name:", person.Name)
	fmt.Println("Age:", person.Age)
	fmt.Println("Profession:", person.Profession)
}
