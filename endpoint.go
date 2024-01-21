package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var employees = []Employee{
	{ID: 1, Name: "John Doe", Role: "Developer"},
	{ID: 2, Name: "Jane Doe", Role: "Designer"},
}

func main() {
	app := fiber.New()

	// Endpoint para obtener todos los empleados
	app.Get("/employees", func(c *fiber.Ctx) error {
		return c.JSON(employees)
	})

	// Endpoint para obtener un empleado por ID
	app.Get("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for _, employee := range employees {
			if fmt.Sprintf("%d", employee.ID) == id {
				return c.JSON(employee)
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
	})

	// Endpoint para agregar un nuevo empleado
	app.Post("/employees", func(c *fiber.Ctx) error {
		var newEmployee Employee
		if err := c.BodyParser(&newEmployee); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Error en el formato de la solicitud"})
		}

		newEmployee.ID = len(employees) + 1
		employees = append(employees, newEmployee)

		return c.JSON(newEmployee)
	})

	// Endpoint para actualizar un empleado por ID
	app.Put("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var updatedEmployee Employee
		if err := c.BodyParser(&updatedEmployee); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Error en el formato de la solicitud"})
		}

		for i, employee := range employees {
			if fmt.Sprintf("%d", employee.ID) == id {
				employees[i] = updatedEmployee
				return c.JSON(updatedEmployee)
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
	})

	// Endpoint para eliminar un empleado por ID
	app.Delete("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, employee := range employees {
			if fmt.Sprintf("%d", employee.ID) == id {
				employees = append(employees[:i], employees[i+1:]...)
				return c.JSON(fiber.Map{"message": "Empleado eliminado exitosamente"})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
	})

	// Iniciar el servidor en el puerto 3000
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
