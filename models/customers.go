package models

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Customer struct {
	CustomerName     NullString
	ContactLastName  NullString
	ContactFirstName NullString
	PhoneNumber      NullString
	AddressLine1     NullString
	AddressLine2     NullString
	City             NullString
	State            NullString
	PostalCode       NullString
	Country          NullString
	CreditLimit      NullFloat64
	Orders           []SimplifiedOrder
	SalesRep         SimplifiedEmployee
}

// GetCustomer will get all  campaigns from a business
func GetCustomer(customerNumber string, c *gin.Context) (Customer, error) {
	var customer Customer
	var simplifiedOrders []SimplifiedOrder
	var orderType Order

	customerSQL := `
	SELECT customerName, C.contactLastName, C.contactFirstName, C.phone, C.addressLine1, C.addressLine2, C.city, C.state, C.postalCode, C.country,
	 C.creditLimit, E.firstName AS salesRepFirstName, E.lastName AS salesRepLastName, E.email AS salesRepEmail, O.orderNumber FROM customers as C 
	 INNER JOIN orders AS O ON C.customerNumber = O.customerNumber
	 INNER JOIN employees AS E ON C.salesRepEmployeeNumber = E.employeeNumber
	 WHERE C.customerNumber=? ;`

	rows, err := db.QueryContext(c, customerSQL, customerNumber)

	if err != nil {
		return customer, err
	}

	var salesRep SimplifiedEmployee
	var order Order
	var orderNumbers []int

	isFirstRow := true

	for rows.Next() {
		if isFirstRow == true {
			if err = rows.Scan(&customer.CustomerName, &customer.ContactLastName, &customer.ContactFirstName,
				&customer.PhoneNumber, &customer.AddressLine1, &customer.AddressLine2, &customer.City, &customer.State, &customer.PostalCode, &customer.Country, &customer.CreditLimit, &salesRep.FirstName, &salesRep.LastName, &salesRep.Email, &order.OrderNumber); err != nil {
				return customer, err
			}
			isFirstRow = false
		} else {
			values := make([]interface{}, 15)
			valuePtr := make([]interface{}, 15)

			for i := 0; i < 15; i++ {
				valuePtr[i] = &values[i]
			}

			if err = rows.Scan(valuePtr...); err != nil {
				return customer, err
			}
			order.OrderNumber = int(values[14].(int64))
		}

		orderNumbers = append(orderNumbers, order.OrderNumber)
	}

	// Get information for each customers
	orderSQL := `SELECT orderNumber, quantityOrdered, priceEach FROM orderdetails WHERE `

	for o := 0; o < len(orderNumbers); o++ {
		if o != 0 {
			orderSQL += ` OR `
		}
		orderSQL += fmt.Sprintf(`orderNumber=%d`, orderNumbers[o])
	}

	orderSQL += ` ORDER BY orderNumber DESC`

	orderRows, err := db.QueryContext(c, orderSQL)

	var previousOrderNumber int
	var orders []Order
	order = orderType

	for orderRows.Next() {
		var orderItem OrderItem

		if err = orderRows.Scan(&orderItem.OrderNumber, &orderItem.QuantityOrdered, &orderItem.PriceEach); err != nil {
			return customer, err
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

	for o := 0; o < len(orders); o++ {
		var simplifiedOrder SimplifiedOrder
		var totalValue float64

		for i := 0; i < len(orders[o].Items); i++ {
			if orders[o].Items[i].PriceEach.Valid == true && orders[o].Items[i].QuantityOrdered.Valid == true {
				totalValue += orders[o].Items[i].PriceEach.Float64 * float64(orders[o].Items[i].QuantityOrdered.Int32)
			}
		}

		simplifiedOrder.OrderNumber = orders[o].OrderNumber
		simplifiedOrder.NumberOfItems = len(orders[o].Items)
		simplifiedOrder.TotalValue = totalValue

		simplifiedOrders = append(simplifiedOrders, simplifiedOrder)
	}

	customer.Orders = simplifiedOrders
	customer.SalesRep = salesRep

	return customer, nil
}
