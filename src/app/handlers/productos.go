package handlers

import (
	"net/http"

	"lab3/src/app/database"
	"lab3/src/app/models"

	"github.com/gin-gonic/gin"
)

func GetProductos(c *gin.Context) {
	var prods []models.Producto
	database.DB.Preload("Categorias").Find(&prods)
	c.JSON(http.StatusOK, prods)
}

func CreateProducto(c *gin.Context) {
	var prod models.Producto
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&prod)
	c.JSON(http.StatusCreated, prod)
}

func UpdateProducto(c *gin.Context) {
	id := c.Param("id")
	var prod models.Producto
	if err := database.DB.First(&prod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no encontrado"})
		return
	}
	if err := c.ShouldBindJSON(&prod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&prod)
	c.JSON(http.StatusOK, prod)
}

func DeleteProducto(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Producto{}, id)
	c.Status(http.StatusNoContent)
}
