package main

import (
	"context"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"inventory/internal/delivery"
	"inventory/internal/repo"
	"inventory/internal/usecase"
)

const (
	jwtSecret    = "your-256-bit-secret" // Замените в production
	databaseName = "inventoryDB"
)

func main() {
	// 1. MongoDB инициализация
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("MongoDB disconnect error: %v", err)
		}
	}()

	// 2. Инициализация репозиториев
	productRepo := repo.NewProductMongoRepo(client, databaseName)
	userRepo := repo.NewUserMongoRepo(client, databaseName)

	// 3. Инициализация use cases
	productUC := usecase.NewProductUseCase(productRepo)
	authUC := usecase.NewAuthUseCase(userRepo)

	// 4. Настройка маршрутизатора
	router := gin.Default()

	// 5. Инициализация обработчиков
	authHandler := delivery.NewAuthRoutes(authUC, []byte(jwtSecret))
	productHandler := delivery.NewProductRoutes(productUC)

	// 6. Регистрация маршрутов
	api := router.Group("/api/v1")
	{
		// Public routes
		authGroup := api.Group("/auth")
		authHandler.RegisterRoutes(authGroup)

		// Protected routes
		authenticated := api.Group("/")
		authenticated.Use(AuthMiddleware())
		productHandler.RegisterRoutes(authenticated)
	}

	// 7. Технические эндпоинты
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "OK",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"database":  "MongoDB",
			"version":   "1.0.0",
		})
	})

	// 8. Запуск сервера
	log.Println("🚀 Server started on :8081")
	if err := router.Run(":8081"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// AuthMiddleware проверка JWT токена
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if username, exists := claims["sub"].(string); exists {
				c.Set("username", username)
			}
		}

		c.Next()
	}
}
