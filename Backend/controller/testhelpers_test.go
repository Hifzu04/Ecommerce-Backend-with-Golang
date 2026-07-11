package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Hifzu04/Ecommerce/Backend/database"
	"github.com/Hifzu04/Ecommerce/Backend/middleware"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// testDB name — isolated from production "Ecommerce" database
const testDBName = "TestEcommerce"

// TestMain runs once before all tests in the controller package.
// It switches collections to a test database and cleans up afterwards.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	// Point global collections to the test database
	testDB := database.Client.Database(testDBName)
	UserCollection = testDB.Collection("Users")
	ProductCollection = testDB.Collection("Products")

	code := m.Run()

	// Cleanup: drop the entire test database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = testDB.Drop(ctx)

	os.Exit(code)
}

// setupRouter creates a Gin router wired exactly like main.go.
func setupRouter() *gin.Engine {
	app := NewApplication(ProductCollection, UserCollection)

	router := gin.New()
	// Public routes (no auth required)
	router.POST("/users/signup", Signup())
	router.POST("/users/login", Login())
	router.POST("/admin/addproducts", Productvieweradmin())
	router.GET("/users/viewproducts", SearchProducts())
	router.GET("/users/searchproducts", SearchProductbyQuery())

	// Protected routes (JWT required)
	protected := router.Group("/")
	protected.Use(middleware.Authentication())
	protected.POST("/addtocart", app.AddToCart())
	protected.GET("/removeitem", app.RemoveItem())
	protected.GET("/listcart", GetItemFromCart())
	protected.POST("/addaddress", AddAddress())
	protected.PUT("/edithomeaddress", EditHomeAddress())
	protected.PUT("/editworkaddress", EditWorkAddress())
	protected.DELETE("/deleteaddresses", DeleteAddress())
	protected.POST("/cartcheckout", app.BuyFromCart())
	protected.GET("/instantbuy", app.InstantBuy())

	return router
}

// cleanCollection removes all documents from a collection.
func cleanCollection(t *testing.T, collName string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	col := database.Client.Database(testDBName).Collection(collName)
	_, _ = col.DeleteMany(ctx, bson.M{})
}

// doRequest is a helper that performs an HTTP request and returns the recorder.
func doRequest(router *gin.Engine, method, url string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, url, reqBody)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
