package handlers

import (
	"net/http"

	"lab3/src/app/database"
	"lab3/src/app/models"

	"github.com/gin-gonic/gin"
)

func GetCategorias(c *gin.Context) {
	var cats []models.Categoria
	database.DB.Preload("Productos").Find(&cats)
	c.JSON(http.StatusOK, cats)
}

func CreateCategoria(c *gin.Context) {
	var cat models.Categoria
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&cat)
	c.JSON(http.StatusCreated, cat)
}

func UpdateCategoria(c *gin.Context) {
	id := c.Param("id")
	var cat models.Categoria
	if err := database.DB.First(&cat, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no encontrada"})
		return
	}
	if err := c.ShouldBindJSON(&cat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&cat)
	c.JSON(http.StatusOK, cat)
}

func DeleteCategoria(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Categoria{}, id)
	c.Status(http.StatusNoContent)
}
