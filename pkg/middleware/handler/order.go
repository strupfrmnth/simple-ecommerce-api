package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/strupfrmnth/simple-ecommerce-api/internal/repository"
)

// Clarify what order handler do
type OrderHandler interface {
	OrderProduct(*gin.Context)
}

type orderHandler struct {
	Repo repository.OrderRepository
}

func NewOrderHandler() OrderHandler {
	return &orderHandler{
		Repo: repository.NewOrderRepository(),
	}
}

func (oh *orderHandler) OrderProduct(c *gin.Context) {
	var orderRequest repository.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// get Float64 will get accurate id
	userID := c.GetFloat64("userID")
	fmt.Println("Ordered by", userID)
	order, err := oh.Repo.CreateOrder(uint(userID), orderRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
