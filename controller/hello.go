package controller

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-jwt-mysql/model"
	"go-jwt-mysql/util"
)

var identityKey = util.GetEnvStr("TOKEN_IDENTITY_KEY", "id")

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims[identityKey],
		"userName": user.(*model.User).Username,
		"text":     "Hello World.",
	})
}
