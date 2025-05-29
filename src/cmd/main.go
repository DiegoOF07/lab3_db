package main

import (
	"fmt"
	"lab3/src/app/database"
	"lab3/src/app/handlers"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const backupFile = "/app/data-out/schema.sql"

func main() {
	DB := database.Connect()

	if err := database.Migrate(DB); err != nil {
		log.Fatalf("Error en migraciones: %v", err)
	}

	ensureBackupFile()

	if err := runPgDump("postgres", "database", "lab3_products", backupFile); err != nil {
		log.Fatalf("Error al ejecutar pg_dump: %v", err)
	}

	info, err := os.Stat(backupFile)
	if err != nil {
		log.Fatalf("No se encontró el archivo %s: %v", backupFile, err)
	}
	if info.Size() == 0 {
		log.Fatalf("El archivo %s se creó pero está vacío", backupFile)
	}

	fmt.Printf("Exportación a %s completada (tamaño: %d bytes)\n", backupFile, info.Size())

	if err := database.SeedData(DB, "/app/data/data.sql"); err != nil {
		log.Fatalf("Error al sembrar datos: %v", err)
	}

	router := setupRouter(DB)
	router.Run(":8080")
}

func ensureBackupFile() {
	if _, err := os.Stat(backupFile); err == nil {
		if err := os.Remove(backupFile); err != nil {
			log.Fatalf("No se pudo eliminar archivo existente %s: %v", backupFile, err)
		}
		fmt.Printf("Se eliminó backup previo: %s\n", backupFile)
	}
}

func runPgDump(user, host, dbname, outFile string) error {
	if err := os.MkdirAll(filepath.Dir(outFile), os.ModePerm); err != nil {
		return fmt.Errorf("no se pudo crear directorio de salida: %w", err)
	}

	cmd := exec.Command(
		"pg_dump",
		"-U", user,
		"-h", host,
		"-d", dbname,
		"-f", outFile,
	)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", os.Getenv("PGPASSWORD")))

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_dump error: %v, salida: %s", err, output)
	}
	return nil
}

func setupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/productos", func(c *gin.Context) {
		handlers.GetProductos(c, db)
	})
	router.GET("/productos/:id", func(c *gin.Context) {
		handlers.GetProducto(c, db)
	})
	router.POST("/productos", func(c *gin.Context) {
		handlers.CreateProducto(c, db)
	})
	router.PUT("/productos/:id", func(c *gin.Context) {
		handlers.UpdateProducto(c, db)
	})
	router.DELETE("/productos/:id", func(c *gin.Context) {
		handlers.DeleteProducto(c, db)
	})
	router.GET("/categorias", func(c *gin.Context) {
		handlers.GetCategorias(c, db)
	})
	return router
}
