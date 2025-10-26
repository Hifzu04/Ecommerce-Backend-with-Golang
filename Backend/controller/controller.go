package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Hifzu04/Ecommerce/Backend/database"
	"github.com/Hifzu04/Ecommerce/Backend/models"
	generate "github.com/Hifzu04/Ecommerce/Backend/tokens"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var UserCollection *mongo.Collection = database.UserData(database.Client, "Users")

var ProductCollection *mongo.Collection = database.ProductData(database.Client, "Products")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)

	}
	return string(bytes)

}

func VerifyPassword(password string, hashpassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password))
	valid := true

	msg := ""

	if err != nil {
		msg = "Pass is incorrect"
		valid = false
	}

	return valid, msg

}

func Signup() gin.HandlerFunc {
	return func(cont *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User

		//c.BindJSON(&user)
		//This tries to parse the incoming HTTP request body as JSON and bind it to the user variable (which is typically a struct).
		if err := cont.BindJSON(&user); err != nil {
			//gin.H is a shortcut for map[string]interface{}.
			//It is used to build key-value pairs for JSON responses easily.
			cont.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := Validate.Struct(user)

		if validationError != nil {
			cont.JSON(http.StatusBadRequest, gin.H{"error": validationError})
			return

		}

		//check if email already exist
		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})

		if err != nil {
			log.Panic(err)
			cont.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			cont.JSON(http.StatusBadRequest, gin.H{"error": "User already exist"})
			return
		}
		//check if phone number already exist
		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			cont.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			cont.JSON(http.StatusBadRequest, gin.H{"error": "phone number already exist"})
			return
		}

		//Hash the pass
		password := HashPassword(*user.Password)
		user.Password = &password

		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		user.ID = primitive.NewObjectID()

		user.User_ID = user.ID.Hex()
		token, refreshtoken, _ := generate.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_ID)
		user.Token = &token
		user.Refresh_Token = &refreshtoken
		user.UserCart = make([]models.ProductUser, 0)
		user.Address_Details = make([]models.Address, 0)
		user.Order_Status = make([]models.Order, 0)
		_, insterterr := UserCollection.InsertOne(ctx, user)

		if insterterr != nil {
			cont.JSON(http.StatusInternalServerError, gin.H{"error": "not created"})
			return
		}
		defer cancel()
		cont.JSON(http.StatusCreated, "Sucessfully signed up!!!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		//C.BindJSON Reads and parses the incoming JSON request body into a Go struct.
		// Used at the beginning of a handler to extract and validate the incoming request data.

		//C.JSON Serializes a Go value (struct, map, slice, etc.) into JSON and writes it to the response.
		//Used to send a JSON response back to the client after processing a request.

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		var dbuser models.User
		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&dbuser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login id is incorrrect"})
			return
		}

		//verify pass
		isValidpass, msg := VerifyPassword(*user.Password, *dbuser.Password)

		if !isValidpass {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}
		token, refreshToken, _ := generate.TokenGenerator(*dbuser.Email, *dbuser.First_Name, *dbuser.Last_Name, dbuser.User_ID)

		generate.UpdateAllTokens(token, refreshToken, user.User_ID)
		c.JSON(http.StatusFound, dbuser)

	}

}

func Productvieweradmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products models.Product
		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		products.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, products)
		if anyerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "not created"})
			return
		}
		c.JSON(http.StatusOK, "Sucessfully added our new product via admin")

	}
}

// to find the product of all the list
func SearchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "error while getting the list of the product")
			return
		}
		//see the documentation
		//https://pkg.go.dev/go.mongodb.org/mongo-driver/v2/mongo#Collection.Find

		if err = cursor.All(ctx, &productlist); err != nil {
			log.Panic(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cancel()
		c.IndentedJSON(200, productlist)

	}
}

// search a/many product(s) by query
func SearchProductbyQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var SearchedPoducts []models.Product
		queryParams := c.Query("name")
		if queryParams == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()
		//$regex: This is a MongoDB operator that stands for regular expression. It allows for powerful pattern-based matching on string fields.
		SearchqueryDB, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": queryParams, "$options": "i"}})

		if err != nil {
			c.IndentedJSON(404, "Something went wrong in fetching the db")
			fmt.Println("error during searching from db %v", err)
		}

		err = SearchqueryDB.All(ctx, &SearchedPoducts)

		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}
		defer SearchqueryDB.Close(ctx)

		if err := SearchqueryDB.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, SearchedPoducts)

	}
}
