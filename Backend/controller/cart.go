package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Hifzu04/Ecommerce/Backend/database"
	"github.com/Hifzu04/Ecommerce/Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	prodcollection *mongo.Collection
	usercollection *mongo.Collection
}

func NewApplication(prodcollection, usercollection *mongo.Collection) *Application {
	return &Application{
		prodcollection: prodcollection,
		usercollection: usercollection,
	}
}

func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
			return
		}

		productId, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)

		defer cancel()

		err = database.addProductToCart(ctx, app.prodcollection, app.usercollection, productId, userQueryID)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "sucessfully added to the cart")

	}

}

func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryId := c.Query("id")
		if productQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}
		userQueryId := c.Query("userID")

		if userQueryId == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("Product id is empty"))
			return
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryId)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 4*time.Second)

		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodcollection, app.usercollection, ProductID, userQueryId)

		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "sucessfully removed from cart")

	}
}

func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			println("id is missing from query")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}
		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var usert models.User

		filter := bson.D{{Key: "_id", Value: usert_id}}
		UserCollection.FindOne(ctx, filter).Decode(&usert)

		//aggregation match -> unwind ->grouping
		// Stage 1: Match documents where _id == usert_id (i.e., only this specific user)
		filter_match := bson.D{
			{Key: "$match", Value: bson.D{{Key: "_id", Value: usert_id}}},
		}

		// Stage 2: Unwind the usercart array → each cart item becomes a separate document
		unwind := bson.D{
			{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$usercart"}}},
		}

		// Stage 3: Group the unwound documents back by _id
		// and calculate the total price of all items in usercart
		grouping := bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$_id"}, // group by user id
				{Key: "total", Value: bson.D{
					{Key: "$sum", Value: "$usercart.price"}, // sum all prices in usercart
				}},
			}},
		}

		// Run the aggregation pipeline: [$match → $unwind → $group]
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})

		if err != nil {
			log.Println(err)
		}

		// Declares a variable listing as a slice of bson.M.
		// bson.M is just a map[string]interface{}, i.e., a flexible JSON-like structure in Go.
		var listing []bson.M
		//pointcursor.All(...) reads all documents from the cursor into the listing slice.
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		for _, json := range listing {
			c.IndentedJSON(200, gin.H{
				"total":    json["total"],
				"usercart": usert.UserCart,
			})
		}
		ctx.Done()

	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")

		if userQueryID == "" {
			log.Panicln("user id is empty ")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("userid is empty"))
		}

		var context, cancel = context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()

		errr := database.BuyItemFromCart(context, app.usercollection, userQueryID)

		if errr != nil {
			c.IndentedJSON(http.StatusInternalServerError, errr)
		}

		c.IndentedJSON(200, "Sucessfully placed the order")

	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("Userid")
		if userQueryID == "" {
			log.Println("User id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("User id is empty"))
		}
		//why pid??
		productQueryid := c.Query("pid")
		if productQueryid == "" {
			log.Println("Product id is empty")
			_ = c.AbortWithError(http.StatusInternalServerError, errors.New("Product id is empty"))
		}

		productID, err := primitive.ObjectIDFromHex(productQueryid)

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)

		defer cancel()
		err = database.InstantBuyer(ctx, app.prodcollection, app.usercollection, productID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Sucessfully placed this order. Thankyou for shopping :) ")
	}

}
