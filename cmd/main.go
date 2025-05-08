package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/juancanchi/products/internal/delivery/http"
	"github.com/juancanchi/products/internal/delivery/http/middleware"
	"github.com/juancanchi/products/internal/infrastructure/postgres"
	"github.com/juancanchi/products/internal/usecase"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=jujuy_market port=5432 sslmode=disable"
	}

	var dbErr error
	var db = tryConnectToDB(dsn, &dbErr)
	if dbErr != nil {
		log.Fatalf("❌ Error al conectar a la base de datos después de varios intentos: %v", dbErr)
	}

	// Inyección de dependencias
	productRepo := postgres.NewProductRepository(db)
	productUC := usecase.NewProductUsecase(productRepo)
	productHandler := http.NewProductHandler(productUC)

	categoryRepo := postgres.NewCategoryRepository(db)
	categoryUC := usecase.NewCategoryUsecase(categoryRepo)
	categoryHandler := http.NewCategoryHandler(categoryUC)

	r := gin.Default()
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecreto"
	}

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"POST", "GET", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Rutas públicas
	r.GET("/products", productHandler.List)
	r.GET("/categories", categoryHandler.List)

	// Rutas autenticadas para usuarios
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware(jwtSecret))
	auth.POST("/products", productHandler.Create)
	auth.GET("/my-products", productHandler.ListByUser)
	auth.GET("/products/:id", productHandler.GetByID)
	auth.PUT("/products/:id", productHandler.Update)
	auth.DELETE("/products/:id", productHandler.Delete)

	// Rutas solo para admin
	admin := r.Group("/")
	admin.Use(middleware.JWTMiddleware(jwtSecret), middleware.AdminOnly())
	admin.PUT("/products/:id/status", productHandler.ChangeStatus)
	admin.POST("/categories", categoryHandler.Create)
	admin.PUT("/categories/:id", categoryHandler.Update)
	admin.DELETE("/categories/:id", categoryHandler.Delete)

	// Puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Servidor escuchando en http://localhost:%s", port)
	r.Run(":" + port)
}

func tryConnectToDB(dsn string, lastErr *error) *gorm.DB {
	const maxAttempts = 10
	for i := 1; i <= maxAttempts; i++ {
		log.Printf("⏳ Intentando conectar a la base de datos... intento %d/%d", i, maxAttempts)
		db, err := postgres.NewDB(dsn)
		if err == nil {
			log.Println("✅ Conexión a la base de datos exitosa.")
			return db
		}
		*lastErr = err
		time.Sleep(3 * time.Second)
	}
	return nil
}
