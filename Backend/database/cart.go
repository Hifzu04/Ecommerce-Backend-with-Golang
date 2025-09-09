package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


func addProductToCart(ctx context.Context , prodCollection , userCollection *mongo.Collection, productID primitive.ObjectID , userID string) error{


}


func RemoveCartItem(ctx context.Context , prodCollection ,userCollection *mongo.Collection , productID primitive.ObjectID , userID string) error{

}


func BuyItemFromCart(ctx context.Context, userCollection *mongo.Collection, userID string) error {

}

func InstantBuyer(ctx context.Context, prodCollection, userCollection *mongo.Collection, productID primitive.ObjectID, UserID string) error {

	
}