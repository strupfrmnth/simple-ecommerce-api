package route

import (
	"log"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	inhandler "github.com/strupfrmnth/simple-ecommerce-api/internal/middleware/handler"
	"github.com/strupfrmnth/simple-ecommerce-api/pkg/middleware/auth"
	"github.com/strupfrmnth/simple-ecommerce-api/pkg/middleware/handler"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
)

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}

func RunAPI(port string, rps int) error {
	limit = ratelimit.New(rps)

	router := gin.Default()
	router.Use(leakBucket())

	userHandler := handler.NewUserHandler()
	productHandler := handler.NewProductHandler()
	orderHandler := handler.NewOrderHandler()
	ratelimitHandler := inhandler.NewRateLimitHandler()

	router.GET("/rate", func(c *gin.Context) {
		c.JSON(200, "rate limiting test")
	})

	router.GET("/", ratelimitHandler.CheckIPLimit, func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to My Ecommerce API "+c.ClientIP())
	})

	userRoute := router.Group("/user")
	{
		userRoute.POST("/register", userHandler.AddUser)
		userRoute.POST("/login", userHandler.LoginUser)
	}

	userProtectedRoute := router.Group("/users", auth.AuthorizeJWT())
	{
		userProtectedRoute.GET("/", userHandler.GetAllUser)
		userProtectedRoute.GET("/:id", userHandler.GetById)
		userProtectedRoute.PUT("/:id", userHandler.UpdateUser)
		userProtectedRoute.DELETE("/:id", userHandler.DeleteUser)
	}

	productRoute := router.Group("/products", auth.AuthorizeJWT())
	{
		productRoute.GET("/", productHandler.GetAllProduct)
		productRoute.GET("/:id", productHandler.GetById)
		productRoute.POST("/", productHandler.AddProduct)
		productRoute.PUT("/:id", productHandler.UpdateProduct)
		productRoute.DELETE("/:id", productHandler.DeleteProduct)
	}

	orderRoute := router.Group("/order", auth.AuthorizeJWT())
	{
		orderRoute.POST("/", orderHandler.OrderProduct)
	}

	return router.Run("localhost:" + port)
}
