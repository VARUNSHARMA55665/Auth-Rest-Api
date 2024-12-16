// routers/router.go
package routers

import (
	"auth-rest-api/middlewares"

	apiController "auth-rest-api/controllers/api"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	//Giving access to storage folder
	r.Static("/storage", "storage")
	r.GET("/swagger/auth-rest-api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Giving access to template folder
	r.Static("/templates", "templates")
	r.LoadHTMLGlob("templates/*")

	userGroup := r.Group("/api/auth-rest-api/user")
	userGroup.Use(middlewares.NoAuthMiddleware())
	{
		userGroup.POST("/signUp", apiController.SignUp)
		userGroup.POST("/signIn", apiController.SignIn)
	}

	return r
}
