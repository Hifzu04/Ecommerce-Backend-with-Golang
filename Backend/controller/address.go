package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/Hifzu04/Ecommerce/Backend/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func AddAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")

		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "Invalid code"})
			c.Abort()
			return
		}
		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, "Internal Server Error")
		}

		var newaddress models.Address
		newaddress.Address_id = primitive.NewObjectID()
		//c.BindJSON(&newaddress) binds the request body JSON into the addresses struct.
		if err = c.BindJSON(&newaddress); err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err)
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		//This aggregation checks how many addresses the user already has:
		//$match – Find the user with _id = user_id.
		//$unwind – Break the address array into individual documents.
		//$group – Group them and count how many addresses exist.
		match_filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

		pointcursor , err := UserCollection.Aggregate(ctx , mongo.Pipeline{match_filter , unwind , group})

		if err != nil {
			c.IndentedJSON(500 , "internal server error")

		}

		var addressinfo []bson.M
		if err = pointcursor.All( ctx ,&addressinfo); err != nil {
			panic(err)
		}
		

	}
}

func EditHomeAddress() {

}

func EditWorkAddress() {

}

func DeleteAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		addresses := make([]models.Address, 0)
		usert_id, err := primitive.ObjectIDFromHex(user_id)

		if err != nil {
			c.IndentedJSON(500, "Internal Server error")
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{bson.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{bson.E{Key: "address", Value: addresses}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(404, "unable to update address")
			return
		}
		defer cancel()

		ctx.Done()

		c.IndentedJSON(200, "sucessfully Deleted the address")

	}

}
