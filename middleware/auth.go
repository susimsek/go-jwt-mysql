package middleware

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-jwt-mysql/model"
	"go-jwt-mysql/payload"
	"go-jwt-mysql/repository"
	"go-jwt-mysql/util"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var tokenExp = util.GetEnvInt("JWT_TOKEN_EXP_SECOND", 86400)
var tokenMaxRefresh = util.GetEnvInt("JWT_TOKEN_MAX_REFRESH_SECOND", 86400)
var tokenSecretKey = util.GetEnvStr("JWT_TOKEN_SECRET", "R1BYcTVXVGNDU2JmWHVnZ1lnN0FKeGR3cU1RUU45QXV4SDJONFZ3ckhwS1N0ZjNCYVkzZ0F4RVBSS1UzRENwRw==")
var identityKey = util.GetEnvStr("TOKEN_IDENTITY_KEY", "id")

func GinJwtMiddlewareHandler() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte(tokenSecretKey),
		Timeout:     time.Duration(int64(time.Second) * int64(tokenExp)),
		MaxRefresh:  time.Duration(int64(time.Second) * int64(tokenMaxRefresh)),
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*model.User); ok {
				return jwt.MapClaims{
					identityKey: v.Username,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &model.User{
				Username: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals payload.Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			var user model.User
			err := repository.GetUserByUsername(&user, userID)
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}

			err = repository.VerifyPassword(user.Password, password)
			if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, jwt.ErrFailedAuthentication
			}

			return &model.User{
				Username: user.Username,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}
