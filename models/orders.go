package models

import (
	"fmt"
	
	"github.com/gin-gonic/gin"
)

// Order is the fully detailled stuct for Order
type Order struct {
	OrderNumber    int
	OrderDate      string
	RequiredDate   string
	ShippedDate    NullString
	Status         string
	Comments       NullString
	CustomerNumber int
	TotalValue     float64
	Items          []OrderItem
}


// SimplifiedOrder is the simplified stuct for Order
type SimplifiedOrder struct {
	OrderNumber   int
	NumberOfItems int
	TotalValue    float64
}

// OrderItem is the fully detailled stuct of and order item
type OrderItem struct {
	OrderNumber     int
	ProductCode     string
	QuantityOrdered int32
	PriceEach       float64
	OrderLineNumber int32
}

// GetOrderDetails will get all Order Item and order
func GetOrderDetails(orderNumber string, c *gin.Context) (Order, error) {
	var order Order
	var orderItems []OrderItem
	var totalValue float64

	orderSQL := `
	SELECT O.orderNumber, O.orderDate, O.requiredDate, O.shippedDate, O.status, O.comments, O.customerNumber FROM orders AS O
	 WHERE O.orderNumber=?;`

	orderRow := db.QueryRow(orderSQL, orderNumber)

	getOrderErr := orderRow.Scan(&order.OrderNumber, &order.OrderDate, &order.RequiredDate, &order.ShippedDate, &order.Status, &order.Comments, &order.CustomerNumber)

	if getOrderErr != nil {
		return order, getOrderErr
	}

	itemsSQL := `
	SELECT O.orderNumber, O.productCode, O.quantityOrdered, O.priceEach, O.orderLineNumber FROM orderdetails as O 
	 WHERE O.orderNumber=?;`

	itemRows, itemsError := db.QueryContext(c, itemsSQL, orderNumber)

	if itemsError != nil {
		return order, itemsError
	}

	for itemRows.Next() {
		var orderItem OrderItem

		if itemErr := itemRows.Scan(&orderItem.OrderNumber, &orderItem.ProductCode, &orderItem.QuantityOrdered, &orderItem.PriceEach, &orderItem.OrderLineNumber); itemErr != nil {
			return order, itemErr
		}

		orderItems = append(orderItems, orderItem)
		totalValue += orderItem.PriceEach * float64(orderItem.QuantityOrdered)
	}

	order.TotalValue = totalValue
	order.Items = orderItems

	return order, nil
}

func getOrdersAndItems(orderNumbers []int, c *gin.Context) ([]Order, error) {
	var orders []Order
	var orderType Order

	orderSQL := `SELECT orderNumber, quantityOrdered, priceEach FROM orderdetails WHERE `

	for o := 0; o < len(orderNumbers); o++ {
		if o != 0 {
			orderSQL += ` OR `
		}
		orderSQL += fmt.Sprintf(`orderNumber=%d`, orderNumbers[o])
	}

	orderSQL += ` ORDER BY orderNumber DESC`

	orderRows, err := db.QueryContext(c, orderSQL)

	if err != nil {
		return orders, err
	}

	var previousOrderNumber int
	var order Order

	for orderRows.Next() {
		var orderItem OrderItem

		if err = orderRows.Scan(&orderItem.OrderNumber, &orderItem.QuantityOrdered, &orderItem.PriceEach); err != nil {
			return orders, err
		}
		if previousOrderNumber == 0 {
			previousOrderNumber = orderItem.OrderNumber
		}
		if previousOrderNumber != orderItem.OrderNumber {
			orders = append(orders, order)
			order = orderType
			order.OrderNumber = orderItem.OrderNumber
		}

		order.Items = append(order.Items, orderItem)
		previousOrderNumber = orderItem.OrderNumber
	}

	return orders, nil
}

func parseSimplifiedOrders(orders []Order) []SimplifiedOrder {
	var simplifiedOrders []SimplifiedOrder

	for o := 0; o < len(orders); o++ {
		var simplifiedOrder SimplifiedOrder
		var totalValue float64

		for i := 0; i < len(orders[o].Items); i++ {
			totalValue += orders[o].Items[i].PriceEach * float64(orders[o].Items[i].QuantityOrdered)
		}

		simplifiedOrder.OrderNumber = orders[o].OrderNumber
		simplifiedOrder.NumberOfItems = len(orders[o].Items)
		simplifiedOrder.TotalValue = totalValue

		simplifiedOrders = append(simplifiedOrders, simplifiedOrder)
	}

	return simplifiedOrders
}