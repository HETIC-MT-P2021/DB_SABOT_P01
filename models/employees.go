package models

import (
	"github.com/gin-gonic/gin"
)

// SimplifiedEmployee is the simplified stuct for Employee
type SimplifiedEmployee struct {
	FirstName string
	LastName  string
	Email     string
}

// Employee is the fully detailled stuct for Employee
type Employee struct {
	EmployeeNumber int
	FirstName string
	LastName  string
	Extension string
	Email     string
	Office Office
	ReportsTo NullInt32
	JobTitle string
}

// GetEmployees will get all employees and their detailled Office
func GetEmployees(c *gin.Context) ([]Employee, error) {
	var employees []Employee

	employeesSQL := `
	 SELECT E.employeeNumber, E.firstName, E.lastName, E.extension, E.email, E.officeCode, E.reportsTo, E.jobTitle, O.city, O.phone, O.addressLine1, O.addressLine2, O.state, O.country, O.postalCode, O.territory   FROM employees AS E
	 INNER JOIN offices AS O ON E.officeCode = O.officeCode;`

	employeeRows, employeesError := db.QueryContext(c, employeesSQL)

	if employeesError != nil {
		return employees, employeesError
	}

	for employeeRows.Next() {
		var employee Employee
		var office Office

		if employeeErr := employeeRows.Scan(&employee.EmployeeNumber, &employee.FirstName, &employee.LastName, &employee.Extension, &employee.Email, &office.OfficeNumber, &employee.ReportsTo, &employee.JobTitle, &office.City, &office.Phone, &office.AddressLine1, &office.AddressLine2, &office.State, &office.Country, &office.PostalCode, &office.Territory); employeeErr != nil {
			return employees, employeeErr
		}

		employee.Office = office

		employees = append(employees, employee)
	}

	return employees, nil
}
