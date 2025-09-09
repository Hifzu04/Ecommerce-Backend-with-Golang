package main

import (
	"log"
	"os"

	"github.com/Hifzu04/Ecommerce/Backend/routes"
	"github.com/gin-gonic/gin"
	"honnef.co/go/tools/config"
)


func main(){
	port:=os.Getenv("PORT")
	if port == ""{
		port= "8000"
	}

	config.ConnectDB()
	router := gin.New()

	router.Use(gin.Logger())
	

	routes.UserRoutes(router)

	
	log.Fatal(router.Run(":" +port))	


}