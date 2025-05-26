CREATE TYPE estado_producto AS ENUM ('disponible', 'agotado', 'discontinuado');

CREATE DOMAIN nombre_categoria AS VARCHAR(50)
    CHECK (VALUE ~* '^[a-zA-Z0-9 ]+$');

CREATE TABLE categorias (
    id SERIAL PRIMARY KEY,
    nombre nombre_categoria NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE productos (
    id SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    precio DECIMAL(10,2) NOT NULL,
    descripcion TEXT,
    estado estado_producto NOT NULL DEFAULT 'disponible',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE productos_categorias (
    producto_id  INTEGER NOT NULL,
    categoria_id INTEGER NOT NULL,
    PRIMARY KEY (producto_id, categoria_id),
    CONSTRAINT fk_producto FOREIGN KEY (producto_id) REFERENCES productos(id) ON DELETE CASCADE,
    CONSTRAINT fk_categoria FOREIGN KEY (categoria_id) REFERENCES categorias(id) ON DELETE CASCADE
);

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

