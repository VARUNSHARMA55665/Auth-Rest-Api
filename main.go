package main

import (
	controllers "auth-rest-api/controllers/api"
	"auth-rest-api/db"
	_ "auth-rest-api/docs"
	"auth-rest-api/helpers"
	"auth-rest-api/models"
	"auth-rest-api/resources"
	"auth-rest-api/routers"
	services "auth-rest-api/services/api"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	r := routers.SetupRouter()
	port := os.Getenv("port")

	resources.GetConfig().AddConfigPath("./resources")
	resources.Start()

	err := db.InitMongoClient()
	if err != nil {
		log.Print("main Error in initiating mongo err: ", err)
		return
	}

	err = helpers.InitRedis()
	if err != nil {
		log.Println("Failed to initialize Redis:", err.Error())
	}

	userProvider := BuildUserProvider()
	controllers.InitUserProvider(userProvider)

	if port == "" {
		port = "8080" //localhost
	}
	r.Run(":" + port)

}

func BuildUserProvider() models.UserProvider {
	return services.InitUser()
}
