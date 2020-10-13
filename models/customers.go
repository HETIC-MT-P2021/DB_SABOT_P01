package models

import (
	"github.com/gin-gonic/gin"
)

// Customer is the fully detailled stuct for Customer
type Customer struct {
	CustomerName     string
	ContactLastName  string
	ContactFirstName string
	PhoneNumber      string
	AddressLine1     string
	AddressLine2     NullString
	City             string
	State            NullString
	PostalCode       NullString
	Country          string
	CreditLimit      NullFloat64
	Orders           []SimplifiedOrder
	SalesRep         SimplifiedEmployee
}

// GetCustomer will get a customer and it's orders
func GetCustomer(customerNumber string, c *gin.Context) (Customer, error) {
	var customer Customer
	var simplifiedOrders []SimplifiedOrder

	customer, orderNumbers, salesRep, getCustErr := getCustomer(customerNumber, c)

	if getCustErr != nil {
		return customer, getCustErr
	}

	orders, getOrderItemErr := getOrdersAndItems(orderNumbers, c)

	if getOrderItemErr != nil {
		return customer, getOrderItemErr
	}

	simplifiedOrders = parseSimplifiedOrders(orders)

	customer.Orders = simplifiedOrders
	customer.SalesRep = salesRep

	return customer, nil
}

func getCustomer(customerNumber string, c *gin.Context ) (Customer, []int, SimplifiedEmployee, error) {
	var customer Customer
	var salesRep SimplifiedEmployee
	var orderNumbers []int

	customerSQL := `
	SELECT customerName, C.contactLastName, C.contactFirstName, C.phone, C.addressLine1, C.addressLine2, C.city, C.state, C.postalCode, C.country,
	 C.creditLimit, E.firstName AS salesRepFirstName, E.lastName AS salesRepLastName, E.email AS salesRepEmail, O.orderNumber FROM customers as C 
	 INNER JOIN orders AS O ON C.customerNumber = O.customerNumber
	 INNER JOIN employees AS E ON C.salesRepEmployeeNumber = E.employeeNumber
	 WHERE C.customerNumber=? ;`

	rows, err := db.QueryContext(c, customerSQL, customerNumber)

	if err != nil {
		return customer, orderNumbers, salesRep, err
	}
	
	var order Order

	isFirstRow := true

	for rows.Next() {
		if isFirstRow == true {
			if err = rows.Scan(&customer.CustomerName, &customer.ContactLastName, &customer.ContactFirstName,
				&customer.PhoneNumber, &customer.AddressLine1, &customer.AddressLine2, &customer.City, &customer.State, &customer.PostalCode, &customer.Country, &customer.CreditLimit, &salesRep.FirstName, &salesRep.LastName, &salesRep.Email, &order.OrderNumber); err != nil {
				return customer, orderNumbers, salesRep, err
			}
			isFirstRow = false
		} else {
			values := make([]interface{}, 15)
			valuePtr := make([]interface{}, 15)

			for i := 0; i < 15; i++ {
				valuePtr[i] = &values[i]
			}

			if err = rows.Scan(valuePtr...); err != nil {
				return customer, orderNumbers, salesRep, err
			}
			order.OrderNumber = int(values[14].(int64))
		}

		orderNumbers = append(orderNumbers, order.OrderNumber)
	}

	return customer, orderNumbers, salesRep, nil
}