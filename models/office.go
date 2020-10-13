package models

import (
	"github.com/gin-gonic/gin"
)

// Office is the fully detailled stuct for office
type Office struct {
	OfficeNumber string
	City string
	Phone  string
	AddressLine1 string
	AddressLine2     NullString
	State NullString
	Country string
	PostalCode string
	Territory string
}

// GetOffice will get an office and it's employee
func GetOffice(officeNumber string, c *gin.Context) (Office, []SimplifiedEmployee, error) {
	var office Office
	var simplifiedEmployees []SimplifiedEmployee

	officeSQL := `
	SELECT O.officeCode, O.city, O.phone, O.addressLine1, O.addressLine2, O.state, O.country, O.postalCode, O.territory, E.firstName, E.lastName, E.email
	 FROM offices AS O
	 INNER JOIN employees AS E ON O.officeCode = E.officeCode WHERE O.officeCode = ?;`

	rows, err := db.QueryContext(c, officeSQL, officeNumber)

	if err != nil {
		return office, simplifiedEmployees, err
	}

	isFirstRow := true

	for rows.Next() {
		var simplifiedEmployee SimplifiedEmployee

		if isFirstRow == true {
			if err = rows.Scan(&office.OfficeNumber, &office.City, &office.Phone, &office.AddressLine1, &office.AddressLine2, &office.State, &office.Country, &office.PostalCode, &office.Territory, &simplifiedEmployee.FirstName, &simplifiedEmployee.LastName, &simplifiedEmployee.Email); err != nil {
				return office, simplifiedEmployees, err
			}
			isFirstRow = false
		} else {
			values := make([]interface{}, 12)
			valuePtr := make([]interface{}, 12)

			for i := 0; i < 12; i++ {
				valuePtr[i] = &values[i]
			}

			if err = rows.Scan(valuePtr...); err != nil {
				return office, simplifiedEmployees, err
			}

			simplifiedEmployee.FirstName = string(values[9].([]uint8))
			simplifiedEmployee.LastName = string(values[10].([]uint8))
			simplifiedEmployee.Email = string(values[11].([]uint8))
		}

		simplifiedEmployees = append(simplifiedEmployees, simplifiedEmployee)
	}

	return office, simplifiedEmployees, nil
}