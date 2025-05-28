package models

import (
	"time"

	"gorm.io/gorm"
)

type EstadoProducto string

const (
	Disponible    EstadoProducto = "disponible"
	Agotado       EstadoProducto = "agotado"
	Discontinuado EstadoProducto = "discontinuado"
)

type Categoria struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nombre    string         `gorm:"type:nombre_categoria;not null;unique" json:"nombre"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Productos []*Producto    `gorm:"many2many:productos_categorias" json:"productos,omitempty"`
}

func (Categoria) TableName() string {
	return "categorias"
}

type Producto struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Nombre      string         `gorm:"size:100;not null" json:"nombre"`
	Precio      float64        `gorm:"type:decimal(10,2);not null" json:"precio"`
	Descripcion string         `json:"descripcion"`
	Estado      EstadoProducto `gorm:"type:estado_producto;not null;default:'disponible'" json:"estado"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Categorias  []*Categoria   `gorm:"many2many:productos_categorias" json:"categorias,omitempty"`
}

func (Producto) TableName() string {
	return "productos"
}

type ProductoCategoriaView struct {
	ProductoID          uint           `gorm:"column:producto_id" json:"producto_id"`
	ProductoNombre      string         `gorm:"column:producto_nombre" json:"producto_nombre"`
	Precio              float64        `gorm:"column:precio" json:"precio"`
	ProductoDescripcion string         `gorm:"column:producto_descripcion" json:"producto_descripcion"`
	ProductoEstado      EstadoProducto `gorm:"column:producto_estado" json:"producto_estado"`
	CategoriaID         uint           `gorm:"column:categoria_id" json:"categoria_id"`
	CategoriaNombre     string         `gorm:"column:categoria_nombre" json:"categoria_nombre"`
}
