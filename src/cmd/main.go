package main

import (
	"lab3/src/app/database"

	// "lab3/src/app/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	// "lab3/src/app/handlers"
)

func main() {
	DB := database.Connect()

	database.Migrate(DB)
	database.SeedData(DB, "/app/data/data.sql")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Rutas Categorías
	// r.GET("/categorias", handlers.GetCategorias)
	// r.POST("/categorias", handlers.CreateCategoria)
	// r.PUT("/categorias/:id", handlers.UpdateCategoria)
	// r.DELETE("/categorias/:id", handlers.DeleteCategoria)

	// Rutas Productos
	// r.GET("/productos", handlers.GetProductos)
	// r.POST("/productos", handlers.CreateProducto)
	// r.PUT("/productos/:id", handlers.UpdateProducto)
	// r.DELETE("/productos/:id", handlers.DeleteProducto)

	router.Run(":8080")
}
