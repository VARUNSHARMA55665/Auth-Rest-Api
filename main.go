package main

import (
	controllers "auth-rest-api/controllers/api"
	"auth-rest-api/db"
	"auth-rest-api/models"
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

	err := db.InitMongoClient()
	if err != nil {
		log.Print("main Error in initiating mongo err: ", err)
		return
	}

	r := routers.SetupRouter()

	userProvider := BuildUserProvider()
	controllers.InitUserProvider(userProvider)

	port := os.Getenv("port")
	if port == "" {
		port = "8080" //localhost
	}
	r.Run(":" + port)

}

func BuildUserProvider() models.UserProvider {
	return services.InitUser()
}
