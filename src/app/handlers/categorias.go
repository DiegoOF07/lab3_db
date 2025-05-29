
package handlers

import (
	"net/http"

	"lab3/src/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetCategorias(c *gin.Context, db *gorm.DB) {
    var categorias []models.Categoria
    if err := db.Table("vista_categorias_completa").Find(&categorias).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener categorías"})
        return
    }
    c.JSON(http.StatusOK, categorias)
}