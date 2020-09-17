package models

type Order struct {
	OrderNumber    int
	OrderDate      string
	RequiredDate   string
	ShippedDate    string
	Status         string
	Comments       string
	CustomerNumber int
	Items          []OrderItem
}

type SimplifiedOrder struct {
	OrderNumber   int
	NumberOfItems int
	TotalValue    float64
}

type OrderItem struct {
	OrderNumber     int
	ProductCode     NullString
	QuantityOrdered NullInt32
	PriceEach       NullFloat64
	OrderLineNumber NullInt32
}
