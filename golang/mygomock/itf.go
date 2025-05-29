package mygomock

type Foo interface {
	Bar(x int) int
}

type Car struct {
	Age   int
	Color string
}

type Dealer interface {
	Evaluate(c *Car) int
}
