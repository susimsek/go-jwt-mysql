package route

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-jwt-mysql/controller"
	"go-jwt-mysql/middleware"
	"log"
)

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddlewareHandler())

	authMiddleware, err := middleware.GinJwtMiddlewareHandler()

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.POST("/register", controller.CreateUser)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", controller.HelloHandler)
		auth.GET("user/:id", controller.GetUserByID)
		auth.PUT("user/:id", controller.UpdateUser)
		auth.DELETE("user/:id", controller.DeleteUser)
		auth.GET("user", controller.GetUsers)
	}

	return r
}
