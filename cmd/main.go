package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	handler "github.com/juancanchi/jujuy-market/products/internal/delivery/http"
	"github.com/juancanchi/jujuy-market/products/internal/infrastructure/postgres"
	"github.com/juancanchi/jujuy-market/products/internal/usecase"
)

func main() {
	// Conexión a la base de datos
	dsn := "host=localhost user=postgres password=postgres dbname=jujuy_market port=5432 sslmode=disable"
	db, err := postgres.NewDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Ejecutar migraciones
	postgres.RunMigrations(db)

	// Inyección de dependencias
	repo := postgres.NewProductRepository(db)
	usecase := usecase.NewProductUsecase(repo)
	handler := handler.NewProductHandler(usecase)

	// Inicializar router
	r := gin.Default()

	// Rutas del microservicio
	r.POST("/products", handler.Create)
	r.GET("/products", handler.List)

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor escuchando en http://localhost:%s", port)
	r.Run(":" + port)
}
