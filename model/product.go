package model

type Product struct {
	ID          int
	Name        string
	Description string
	Price       float64
}

type AllProducts struct {
	Products []*Product
}
