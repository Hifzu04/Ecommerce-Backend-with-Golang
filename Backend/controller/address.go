package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Hifzu04/Ecommerce/Backend/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
		defer cancel()

		//This aggregation checks how many addresses the user already has:
		//$match – Find the user with _id = user_id.
		//$unwind – Break the address array into individual documents.
		//$group – Group them and count how many addresses exist.
		match_filter := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$address"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$address_id"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}}

		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})

		if err != nil {
			c.IndentedJSON(500, "internal server error")

		}

		var addressinfo []bson.M
		if err = pointcursor.All(ctx, &addressinfo); err != nil {
			panic(err)
		}
		var size int32
		for _, address_no := range addressinfo {
			count := address_no["count"]
			size = count.(int32)
		}
		// maximum 2 address is allowed for a particular user

		if size < 2 {
			filter := bson.D{{Key: "_id", Value: usert_id}}
			update := bson.D{{Key: "$push", Value: bson.D{{Key: "address", Value: newaddress}}}}
			_, err := UserCollection.UpdateOne(ctx, filter, update)
			if err != nil {
				fmt.Println(err)
			} else {
				c.IndentedJSON(400, "not allowed")
			}
		}
		defer cancel()
		ctx.Done()

	}
}

func EditHomeAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		userid := c.Query("id")
		if userid == "" {
			c.Header("Content-Type", "Application//json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editAddress models.Address
		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "address.0.house_name", Value: editAddress.House},
			{Key: "address.0.street_name", Value: editAddress.Street},
			{Key: "address.0.city_name", Value: editAddress.City}, {Key: "address.0.pin_code", Value: editAddress.Pincode}}}}
		_, err = UserCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "something went wrong")
			return

		}

		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "Sucessfully updated the work address")

	}
}

func EditWorkAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Wrong id not provided"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(500, err)
		}
		var editAddress models.Address

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		if err := c.BindJSON(&editAddress); err != nil {
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}
		filter := bson.D{{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "address.1.house_name", Value: editAddress.House},
			{Key: "address.1.street_name", Value: editAddress.Street}, {Key: "address.1.city_name", Value: editAddress.City},
			{Key: "address.1.pin_code", Value: editAddress.Pincode}}}}

		_, err = UserCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			c.IndentedJSON(404, "Wrong")
			return
		}

		defer cancel()

		ctx.Done()
		c.IndentedJSON(200, "Sucessfully deleted")

	}
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
