package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/strupfrmnth/simple-ecommerce-api/internal/repository"
	"github.com/strupfrmnth/simple-ecommerce-api/pkg/middleware/auth"
)

// Clarify what user handler do
type UserHandler interface {
	LoginUser(*gin.Context)
	AddUser(*gin.Context)
	GetAllUser(*gin.Context)
	GetById(*gin.Context)
	UpdateUser(*gin.Context)
	DeleteUser(*gin.Context)
}

type userHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler() UserHandler {
	return &userHandler{
		Repo: repository.NewUserRepository(),
	}
}

func (uh *userHandler) AddUser(c *gin.Context) {
	var user repository.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Repo.AddUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *userHandler) LoginUser(c *gin.Context) {
	var requser repository.User
	// bind a request body into a type, reference: https://gin-gonic.com/docs/examples/binding-and-validation/
	if err := c.ShouldBindJSON(&requser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get user from database
	user, err := uh.Repo.GetByUsername(requser.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}
	signature := auth.CreateToken(user.ID, user.Name)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in", "token": signature})
}

func (uh *userHandler) GetAllUser(c *gin.Context) {
	users, err := uh.Repo.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (uh *userHandler) GetById(c *gin.Context) {
	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.Repo.GetById(intID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *userHandler) UpdateUser(c *gin.Context) {
	var user repository.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = uint(intID)
	user, err = uh.Repo.UpdateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *userHandler) DeleteUser(c *gin.Context) {
	stringID := c.Param("id")
	intID, err := strconv.Atoi(stringID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := uh.Repo.Delete(intID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Successfully delete user"})
}
