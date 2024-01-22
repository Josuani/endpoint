package main

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

var (
	db *gorm.DB
)

func init() {
	// Conectar a la base de datos SQLite
	var err error
	db, err = gorm.Open(sqlite.Open("employees.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Error de conexión a la base de datos:", err)
	}

	// Migrar el modelo Employee a la base de datos
	db.AutoMigrate(&Employee{})
}

func main() {
	app := fiber.New()

	// Endpoint para obtener todos los empleados
	app.Get("/employees", func(c *fiber.Ctx) error {
		var employees []Employee
		result := db.Find(&employees)
		if result.Error != nil {
			log.Println("Error al obtener empleados:", result.Error)
			return c.Status(500).JSON(fiber.Map{"error": "Error interno del servidor"})
		}
		return c.JSON(employees)
	})

	// Endpoint para obtener un empleado por ID
	app.Get("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var employee Employee
		result := db.First(&employee, id)
		if result.Error != nil {
			log.Println("Error al obtener empleado:", result.Error)
			return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
		}
		return c.JSON(employee)
	})

	// Endpoint para agregar un nuevo empleado
	app.Post("/employees", func(c *fiber.Ctx) error {
		var newEmployee Employee
		if err := c.BodyParser(&newEmployee); err != nil {
			log.Println("Error al parsear el cuerpo de la solicitud:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Error en el formato de la solicitud"})
		}

		result := db.Create(&newEmployee)
		if result.Error != nil {
			log.Println("Error al crear un nuevo empleado:", result.Error)
			return c.Status(500).JSON(fiber.Map{"error": "Error interno del servidor"})
		}

		return c.JSON(newEmployee)
	})

	// Endpoint para actualizar un empleado por ID
	app.Put("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		var updatedEmployee Employee
		if err := c.BodyParser(&updatedEmployee); err != nil {
			log.Println("Error al parsear el cuerpo de la solicitud:", err)
			return c.Status(400).JSON(fiber.Map{"error": "Error en el formato de la solicitud"})
		}

		result := db.First(&Employee{}, id)
		if result.Error != nil {
			log.Println("Error al obtener empleado para actualizar:", result.Error)
			return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
		}

		result = db.Model(&Employee{}).Where("id = ?", id).Updates(updatedEmployee)
		if result.Error != nil {
			log.Println("Error al actualizar empleado:", result.Error)
			return c.Status(500).JSON(fiber.Map{"error": "Error interno del servidor"})
		}

		return c.JSON(updatedEmployee)
	})

	// Endpoint para eliminar un empleado por ID
	app.Delete("/employees/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		result := db.Delete(&Employee{}, id)
		if result.Error != nil || result.RowsAffected == 0 {
			log.Println("Error al eliminar empleado:", result.Error)
			return c.Status(404).JSON(fiber.Map{"error": "Empleado no encontrado"})
		}

		return c.JSON(fiber.Map{"message": "Empleado eliminado exitosamente"})
	})

	// Iniciar el servidor en el puerto 3000
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}
