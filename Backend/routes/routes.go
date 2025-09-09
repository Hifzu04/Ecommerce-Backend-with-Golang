package routes

import "github.com/gin-gonic/gin"

var app = controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

func UserRoutes(router *gin.Engine) {
	router.POST("users/signup", controller.Signup())
	router.POST("users/login", controller.Login())
	router.POST("/admin/addproducts", controller.Productvieweradmin())
	router.GET("users/viewproducts", controller.SearchProducts())
	router.GET("users/searchproducts", controller.SearchProductbyQuery())


	router.Use(middleware.Authentication())
	{
		router.GET("/addtocart", app.Addtocart())
		router.GET("removeitem")
		router.GET("/checkout")
		router.GET("/instantbuy")
	}
}



// package routes
 
// import (
//     "github.com/gin-gonic/gin"
//     controller "github.com/golangcompany/JWT-Authentication/controllers"
//     "github.com/golangcompany/JWT-Authentication/middleware"
// )
 
// func UserRoutes(incomingRoutes *gin.Engine) {
//     incomingRoutes.POST("users/signup", controller.Signup())
//     incomingRoutes.POST("users/login", controller.Login())
// }
// func AuthRoutes(incomingRoutes *gin.Engine) {
//     incomingRoutes.Use(middleware.UserAuthenticate())
//     incomingRoutes.GET("/usersdata", controller.GetUsers())
//     incomingRoutes.GET("/users/:user_id", controller.GetUser())
// }