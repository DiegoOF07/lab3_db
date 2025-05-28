package database

import (
	"lab3/src/app/models"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'estado_producto') THEN
                CREATE TYPE estado_producto AS ENUM ('disponible','agotado','discontinuado');
            END IF;
        END$$;

        DO $$
        BEGIN
            CREATE DOMAIN nombre_categoria AS VARCHAR(50) CHECK (VALUE ~* '^[a-zA-Z0-9 \u00C0-\u00FF]+$');
        END$$;`).Error; err != nil {
		return err
	}

	if err := db.AutoMigrate(
		&models.Categoria{},
		&models.Producto{},
	); err != nil {
		return err
	}

	if err := db.Exec(`
        CREATE OR REPLACE VIEW vista_producto_categoria AS
		SELECT
			p.id AS producto_id,
			p.nombre AS producto_nombre,
			p.precio,
			p.descripcion AS producto_descripcion,
			p.estado AS producto_estado,
			c.id AS categoria_id,
			c.nombre AS categoria_nombre
		FROM productos p
		JOIN productos_categorias pc ON p.id = pc.producto_id
		JOIN categorias c ON c.id = pc.categoria_id
		WHERE p.deleted_at IS NULL
		AND c.deleted_at IS NULL;
    `).Error; err != nil {
		return err
	}

	return nil
}
