package main

import (
	"github.com/gin-contrib/cors"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	handler "github.com/juancanchi/products/internal/delivery/http"
	"github.com/juancanchi/products/internal/delivery/http/middleware"
	"github.com/juancanchi/products/internal/infrastructure/postgres"
	"github.com/juancanchi/products/internal/usecase"
)

func main() {
	// Conexión a la base de datos
	dsn := "host=localhost user=postgres password=postgres dbname=jujuy_market port=5432 sslmode=disable"
	db, err := postgres.NewDB(dsn)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}

	// Inyección de dependencias
	repo := postgres.NewProductRepository(db)
	usecase := usecase.NewProductUsecase(repo)
	handler := handler.NewProductHandler(usecase)

	// Inicializar router
	r := gin.Default()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecreto"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware(jwtSecret))
	auth.POST("/products", handler.Create)
	auth.GET("/my-products", handler.ListByUser)
	auth.GET("/products/:id", handler.GetByID)
	auth.PUT("/products/:id", handler.Update)
	auth.DELETE("/products/:id", handler.Delete)

	r.GET("/products", handler.List)

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor escuchando en http://localhost:%s", port)
	r.Run(":" + port)
}
