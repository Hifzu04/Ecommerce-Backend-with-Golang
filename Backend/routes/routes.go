package routes

import (
	"github.com/Hifzu04/Ecommerce/Backend/controller"
	"github.com/gin-gonic/gin"
)



func UserRoutes(router *gin.Engine) {
	router.POST("users/signup", controller.Signup())
	router.POST("users/login", controller.Login())
	router.POST("/admin/addproducts", controller.Productvieweradmin())
	router.GET("users/viewproducts", controller.SearchProducts())
	router.GET("users/searchproducts", controller.SearchProductbyQuery())

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
