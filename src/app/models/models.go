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
	ID        uint   `gorm:"primaryKey"`
	Nombre    string `gorm:"type:nombre_categoria;not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Productos []*Producto    `gorm:"many2many:productos_categorias"`
}

func (Categoria) TableName() string {
	return "categorias"
}

type Producto struct {
	ID          uint    `gorm:"primaryKey"`
	Nombre      string  `gorm:"size:100;not null"`
	Precio      float64 `gorm:"type:decimal(10,2);not null"`
	Descripcion string
	Estado      EstadoProducto `gorm:"type:estado_producto;not null;default:'disponible'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Categorias  []*Categoria   `gorm:"many2many:productos_categorias"`
}

func (Producto) TableName() string {
	return "productos"
}
