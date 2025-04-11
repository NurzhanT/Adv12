package delivery

import (
	"net/http"
	"time"

	"inventory/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthRoutes struct {
	authUC    usecase.AuthUseCase
	jwtSecret []byte
}

func NewAuthRoutes(authUC usecase.AuthUseCase, jwtSecret []byte) *AuthRoutes {
	return &AuthRoutes{
		authUC:    authUC,
		jwtSecret: jwtSecret,
	}
}

func (h *AuthRoutes) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

func (h *AuthRoutes) register(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required,min=3"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authUC.Register(request.Username, request.Password); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (h *AuthRoutes) login(c *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUC.Login(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.Username,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
		"role": "user",
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": tokenString,
		"token_type":   "Bearer",
		"expires_in":   int64((time.Hour * 24).Seconds()),
	})
}
