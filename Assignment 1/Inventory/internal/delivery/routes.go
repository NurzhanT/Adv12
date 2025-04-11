package delivery

import (
	"inventory/internal/domain"
	"inventory/internal/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterProductRoutes(r gin.IRoutes, uc usecase.ProductUseCase) {
	r.POST("/products", func(c *gin.Context) {
		var p domain.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		uc.Create(p)
		c.JSON(http.StatusCreated, p)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		p, err := uc.GetByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	r.PATCH("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var p domain.Product
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err := uc.Update(id, p)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, p)
	})

	r.DELETE("/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		err := uc.Delete(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	r.GET("/products", func(c *gin.Context) {
		list, _ := uc.List()
		c.JSON(http.StatusOK, list)
	})
}
