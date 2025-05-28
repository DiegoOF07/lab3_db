package handlers

import (
	"net/http"
	"time"

	"lab3/src/app/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetProductos(c *gin.Context, db *gorm.DB) {
	var products []models.ProductoCategoriaView
	if err := db.Table("vista_producto_categoria").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener productos"})
		return
	}

	productosMap := make(map[uint]*models.Producto)
	for _, pv := range products {
		if producto, exists := productosMap[pv.ProductoID]; exists {
			categoria := &models.Categoria{
				ID:     pv.CategoriaID,
				Nombre: pv.CategoriaNombre,
			}
			producto.Categorias = append(producto.Categorias, categoria)
		} else {
			producto := &models.Producto{
				ID:          pv.ProductoID,
				Nombre:      pv.ProductoNombre,
				Precio:      pv.Precio,
				Descripcion: pv.ProductoDescripcion,
				Estado:      pv.ProductoEstado,
				Categorias: []*models.Categoria{
					{
						ID:     pv.CategoriaID,
						Nombre: pv.CategoriaNombre,
					},
				},
			}
			productosMap[pv.ProductoID] = producto
		}
	}

	var productos []*models.Producto
	for _, p := range productosMap {
		productos = append(productos, p)
	}

	c.JSON(http.StatusOK, productos)
}

func GetProducto(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var productosView []models.ProductoCategoriaView
	if err := db.Table("vista_producto_categoria").Where("producto_id = ?", id).Find(&productosView).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto"})
		return
	}

	if len(productosView) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	producto := models.Producto{
		ID:          productosView[0].ProductoID,
		Nombre:      productosView[0].ProductoNombre,
		Precio:      productosView[0].Precio,
		Descripcion: productosView[0].ProductoDescripcion,
		Estado:      productosView[0].ProductoEstado,
	}

	for _, pv := range productosView {
		producto.Categorias = append(producto.Categorias, &models.Categoria{
			ID:     pv.CategoriaID,
			Nombre: pv.CategoriaNombre,
		})
	}

	c.JSON(http.StatusOK, producto)
}

func CreateProducto(c *gin.Context, db *gorm.DB) {
	var productoInput struct {
		Nombre       string  `json:"nombre" binding:"required"`
		Precio       float64 `json:"precio" binding:"required"`
		Descripcion  string  `json:"descripcion"`
		Estado       string  `json:"estado"`
		CategoriaIDs []uint  `json:"categoria_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&productoInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estado := models.EstadoProducto(productoInput.Estado)
	if estado != models.Disponible && estado != models.Agotado && estado != models.Discontinuado {
		estado = models.Disponible
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	producto := models.Producto{
		Nombre:      productoInput.Nombre,
		Precio:      productoInput.Precio,
		Descripcion: productoInput.Descripcion,
		Estado:      estado,
	}

	if err := tx.Create(&producto).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al crear producto"})
		return
	}

	for _, catID := range productoInput.CategoriaIDs {
		var categoria models.Categoria
		if err := tx.First(&categoria, catID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Categoría no encontrada"})
			return
		}

		if err := tx.Exec("INSERT INTO productos_categorias (producto_id, categoria_id) VALUES (?, ?)", producto.ID, catID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asociar categoría"})
			return
		}
	}

	tx.Commit()

	var productoView []models.ProductoCategoriaView
	if err := db.Table("vista_producto_categoria").Where("producto_id = ?", producto.ID).Find(&productoView).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto creado"})
		return
	}

	responseProducto := models.Producto{
		ID:          productoView[0].ProductoID,
		Nombre:      productoView[0].ProductoNombre,
		Precio:      productoView[0].Precio,
		Descripcion: productoView[0].ProductoDescripcion,
		Estado:      productoView[0].ProductoEstado,
	}

	for _, pv := range productoView {
		responseProducto.Categorias = append(responseProducto.Categorias, &models.Categoria{
			ID:     pv.CategoriaID,
			Nombre: pv.CategoriaNombre,
		})
	}

	c.JSON(http.StatusCreated, responseProducto)
}

func UpdateProducto(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	var productoInput struct {
		Nombre      string   `json:"nombre"`
		Precio      float64  `json:"precio"`
		Descripcion string   `json:"descripcion"`
		Estado      string   `json:"estado"`
		CategoriaIDs []uint  `json:"categoria_ids"`
	}

	if err := c.ShouldBindJSON(&productoInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var producto models.Producto
	if err := tx.First(&producto, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	if productoInput.Nombre != "" {
		producto.Nombre = productoInput.Nombre
	}
	if productoInput.Precio != 0 {
		producto.Precio = productoInput.Precio
	}
	if productoInput.Descripcion != "" {
		producto.Descripcion = productoInput.Descripcion
	}
	if productoInput.Estado != "" {
		estado := models.EstadoProducto(productoInput.Estado)
		if estado == models.Disponible || estado == models.Agotado || estado == models.Discontinuado {
			producto.Estado = estado
		}
	}

	producto.UpdatedAt = time.Now()

	if err := tx.Save(&producto).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar producto"})
		return
	}

	if productoInput.CategoriaIDs != nil {
		if err := tx.Exec("DELETE FROM productos_categorias WHERE producto_id = ?", id).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al actualizar categorías"})
			return
		}

		for _, catID := range productoInput.CategoriaIDs {
			var categoria models.Categoria
			if err := tx.First(&categoria, catID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Categoría no encontrada"})
				return
			}

			if err := tx.Exec("INSERT INTO productos_categorias (producto_id, categoria_id) VALUES (?, ?)", id, catID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al asociar categoría"})
				return
			}
		}
	}

	tx.Commit()

	var productoView []models.ProductoCategoriaView
	if err := db.Table("vista_producto_categoria").Where("producto_id = ?", id).Find(&productoView).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener el producto actualizado"})
		return
	}

	responseProducto := models.Producto{
		ID:          productoView[0].ProductoID,
		Nombre:      productoView[0].ProductoNombre,
		Precio:      productoView[0].Precio,
		Descripcion: productoView[0].ProductoDescripcion,
		Estado:      productoView[0].ProductoEstado,
	}

	for _, pv := range productoView {
		responseProducto.Categorias = append(responseProducto.Categorias, &models.Categoria{
			ID:     pv.CategoriaID,
			Nombre: pv.CategoriaNombre,
		})
	}

	c.JSON(http.StatusOK, responseProducto)
}

func DeleteProducto(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")

	result := db.Delete(&models.Producto{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al eliminar producto"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}
