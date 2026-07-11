package main

import (
	"log"
	"os"

	"github.com/Hifzu04/Ecommerce/Backend/controller"
	"github.com/Hifzu04/Ecommerce/Backend/database"
	_ "github.com/Hifzu04/Ecommerce/Backend/docs" // auto-generated docs
	"github.com/Hifzu04/Ecommerce/Backend/middleware"
	"github.com/Hifzu04/Ecommerce/Backend/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Ecommerce Backend API
// @version         1.0
// @description     REST API for an e-commerce platform with user auth, products, cart, orders, and addresses.
// @host            localhost:8000
// @BasePath        /
// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        token
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	///app := controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	var app = controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	//records details about each HTTP request
	routes.UserRoutes(router)
   //Open http://localhost:8000/swagger/index.html in your browser — you'll see the interactive Swagger UI.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(middleware.Authentication())
	router.POST("/addtocart", app.AddToCart())
	router.GET("/removeitem", app.RemoveItem())
	router.GET("/listcart", controller.GetItemFromCart())
	router.POST("/addaddress", controller.AddAddress())
	router.PUT("/edithomeaddress", controller.EditHomeAddress())
	router.PUT("/editworkaddress", controller.EditWorkAddress())
	router.DELETE("/deleteaddresses", controller.DeleteAddress())
	router.POST("/cartcheckout", app.BuyFromCart())
	router.GET("/instantbuy", app.InstantBuy())
	log.Fatal(router.Run(":" + port))
}
