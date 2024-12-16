// routers/router.go
package routers

import (
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

	// r.Use(auth.Middleware())

	// v1 := r.Group("/api/v1")
	// {
	// 	v1.GET("/delivery", controllers.DeliveryHandler)
	// }

	// verifyOtpGroup := r.Group("/api/dao/v1/kyc/verifyOtp")
	// verifyOtpGroup.Use(middlewares.NoAuthMiddleware())
	// {
	// 	verifyOtpGroup.POST("", apiControllerV1.VerifyOtp)
	// }

	return r
}
