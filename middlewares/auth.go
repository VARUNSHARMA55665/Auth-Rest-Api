package middlewares

import (
	apihelpers "auth-rest-api/apiHelpers"
	"auth-rest-api/constants"
	"auth-rest-api/helpers"
	"auth-rest-api/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NoAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//Code for middlewares
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var resJS apihelpers.APIRes
		var reqH models.ReqHeader
		if err := c.ShouldBindHeader(&reqH); err != nil {
			log.Println("AuthMiddleware error in ShouldBindHeader err: ", err)
			c.JSON(http.StatusInternalServerError, resJS)
			c.Abort()
		}

		if len(reqH.Authorization) <= 7 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.TokenMissing]
			resJS.ErrorCode = constants.TokenMissing
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		sub, err := helpers.ValidateToken(reqH.Authorization[7:])
		if err != nil {
			log.Println("AuthMiddleware Error validating auth token err: ", err)
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		var keyExists int64
		keyExists, _ = helpers.Exists("auth|" + sub).Result()
		if keyExists != 1 {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusForbidden, resJS)
			c.Abort()
			return
		}

		keyVal := helpers.GetRedis("auth|" + sub).Val()
		if !strings.EqualFold(keyVal, reqH.Authorization[7:]) {
			resJS.Status = false
			resJS.Message = constants.ErrorCodeMap[constants.InvalidToken]
			resJS.ErrorCode = constants.InvalidToken
			c.JSON(http.StatusUnauthorized, resJS)
			c.Abort()
			return
		}

		c.Next()
	}
}
