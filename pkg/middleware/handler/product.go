package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/strupfrmnth/simple-ecommerce-api/internal/repository"
)

// Clarify what product handler do
type ProductHandler interface {
	AddProduct(*gin.Context)
	GetAllProduct(*gin.Context)
	GetById(*gin.Context)
	UpdateProduct(*gin.Context)
	DeleteProduct(*gin.Context)
}

type productHandler struct {
	Repo repository.ProductRepository
}

func NewProductHandler() ProductHandler {
	return &productHandler{
		Repo: repository.NewProductRepository(),
	}
}

func (ph *productHandler) AddProduct(c *gin.Context) {
	var product repository.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ph.Repo.AddProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ph *productHandler) GetAllProduct(c *gin.Context) {
	products, err := ph.Repo.GetAllProduct()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (ph *productHandler) GetById(c *gin.Context) {
	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product, err := ph.Repo.GetById(intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ph *productHandler) UpdateProduct(c *gin.Context) {
	var product repository.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	product.ID = uint(intID)
	product, err = ph.Repo.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (ph *productHandler) DeleteProduct(c *gin.Context) {
	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := ph.Repo.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Successfully delete product"})
}
