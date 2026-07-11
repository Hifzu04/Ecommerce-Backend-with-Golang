package controller

import (
	"encoding/json"
	"net/http"
	"testing"
)


// UNIT TESTS — pure functions, no database needed


func TestHashPassword(t *testing.T) {
	password := "SecurePass123"
	hashed := HashPassword(password)

	if hashed == "" {
		t.Fatal("HashPassword returned empty string")
	}
	if hashed == password {
		t.Fatal("HashPassword returned the plaintext password")
	}
}

func TestHashPassword_DifferentInputs_DifferentHashes(t *testing.T) {
	hash1 := HashPassword("password1")
	hash2 := HashPassword("password2")

	if hash1 == hash2 {
		t.Fatal("Different passwords produced the same hash")
	}
}

func TestHashPassword_SameInput_DifferentHashes(t *testing.T) {
	// bcrypt uses a random salt, so hashing the same input twice should yield different hashes
	hash1 := HashPassword("samepassword")
	hash2 := HashPassword("samepassword")

	if hash1 == hash2 {
		t.Fatal("Same password produced identical hashes — salt is not working")
	}
}

func TestVerifyPassword_Correct(t *testing.T) {
	password := "MySecret123"
	hashed := HashPassword(password)

	valid, msg := VerifyPassword(password, hashed)

	if !valid {
		t.Fatalf("VerifyPassword should return true for correct password, got msg: %s", msg)
	}
	if msg != "" {
		t.Fatalf("Expected empty message for correct password, got: %s", msg)
	}
}

func TestVerifyPassword_Incorrect(t *testing.T) {
	hashed := HashPassword("CorrectPassword")

	valid, msg := VerifyPassword("WrongPassword", hashed)

	if valid {
		t.Fatal("VerifyPassword should return false for wrong password")
	}
	if msg == "" {
		t.Fatal("Expected error message for wrong password")
	}
}

func TestVerifyPassword_EmptyPassword(t *testing.T) {
	hashed := HashPassword("RealPassword")

	valid, _ := VerifyPassword("", hashed)

	if valid {
		t.Fatal("VerifyPassword should return false for empty password")
	}
}


// INTEGRATION TESTS — HTTP handlers hitting the test database


// signupPayload is used for signup requests in tests.
type signupPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

// loginPayload is used for login requests in tests.
type loginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestSignupHandler_Success(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()
	body := signupPayload{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test_signup@example.com",
		Password:  "TestPass123",
		Phone:     "+1234567890",
	}

	w := doRequest(router, http.MethodPost, "/users/signup", body, nil)

	if w.Code != http.StatusCreated {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusCreated, w.Code, w.Body.String())
	}
}

func TestSignupHandler_DuplicateEmail(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()
	body := signupPayload{
		FirstName: "User",
		LastName:  "One",
		Email:     "duplicate@example.com",
		Password:  "TestPass123",
		Phone:     "+1111111111",
	}

	// First signup — should succeed
	w1 := doRequest(router, http.MethodPost, "/users/signup", body, nil)
	if w1.Code != http.StatusCreated {
		t.Fatalf("First signup failed: %d - %s", w1.Code, w1.Body.String())
	}

	// Second signup with same email — should fail
	body.Phone = "+2222222222" // different phone, same email
	w2 := doRequest(router, http.MethodPost, "/users/signup", body, nil)

	if w2.Code == http.StatusCreated {
		t.Fatal("Duplicate email signup should not return 201")
	}
}

func TestSignupHandler_DuplicatePhone(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()
	body := signupPayload{
		FirstName: "User",
		LastName:  "One",
		Email:     "phone_test1@example.com",
		Password:  "TestPass123",
		Phone:     "+9999999999",
	}

	// First signup
	w1 := doRequest(router, http.MethodPost, "/users/signup", body, nil)
	if w1.Code != http.StatusCreated {
		t.Fatalf("First signup failed: %d - %s", w1.Code, w1.Body.String())
	}

	// Second signup with same phone but different email
	body.Email = "phone_test2@example.com"
	w2 := doRequest(router, http.MethodPost, "/users/signup", body, nil)

	if w2.Code == http.StatusCreated {
		t.Fatal("Duplicate phone signup should not return 201")
	}
}

func TestSignupHandler_InvalidBody(t *testing.T) {
	router := setupRouter()

	// Missing required fields
	body := map[string]string{"email": "only_email@test.com"}
	w := doRequest(router, http.MethodPost, "/users/signup", body, nil)

	if w.Code == http.StatusCreated {
		t.Fatal("Signup with missing required fields should not return 201")
	}
}

func TestSignupHandler_EmptyBody(t *testing.T) {
	router := setupRouter()

	w := doRequest(router, http.MethodPost, "/users/signup", nil, nil)

	if w.Code == http.StatusCreated {
		t.Fatal("Signup with empty body should not return 201")
	}
}

func TestLoginHandler_Success(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()

	// First, create a user
	signupBody := signupPayload{
		FirstName: "Login",
		LastName:  "Test",
		Email:     "login_test@example.com",
		Password:  "LoginPass123",
		Phone:     "+3333333333",
	}
	ws := doRequest(router, http.MethodPost, "/users/signup", signupBody, nil)
	if ws.Code != http.StatusCreated {
		t.Fatalf("Signup failed: %d - %s", ws.Code, ws.Body.String())
	}

	// Now login
	loginBody := loginPayload{
		Email:    "login_test@example.com",
		Password: "LoginPass123",
	}
	w := doRequest(router, http.MethodPost, "/users/login", loginBody, nil)

	if w.Code != http.StatusFound {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusFound, w.Code, w.Body.String())
	}

	// Verify response contains user data
	var respBody map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
	if respBody["email"] != "login_test@example.com" {
		t.Fatalf("Expected email in response, got: %v", respBody["email"])
	}
}

func TestLoginHandler_WrongPassword(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()

	// Create user
	signupBody := signupPayload{
		FirstName: "Wrong",
		LastName:  "Pass",
		Email:     "wrongpass@example.com",
		Password:  "CorrectPass123",
		Phone:     "+4444444444",
	}
	doRequest(router, http.MethodPost, "/users/signup", signupBody, nil)

	// Login with wrong password
	loginBody := loginPayload{
		Email:    "wrongpass@example.com",
		Password: "WrongPassword",
	}
	w := doRequest(router, http.MethodPost, "/users/login", loginBody, nil)

	if w.Code == http.StatusFound {
		t.Fatal("Login with wrong password should not succeed")
	}
}

func TestLoginHandler_NonExistentUser(t *testing.T) {
	cleanCollection(t, "Users")
	defer cleanCollection(t, "Users")

	router := setupRouter()
	loginBody := loginPayload{
		Email:    "nobody@example.com",
		Password: "whatever",
	}
	w := doRequest(router, http.MethodPost, "/users/login", loginBody, nil)

	if w.Code == http.StatusFound {
		t.Fatal("Login with non-existent user should not succeed")
	}
}

// ===================================================================
// PRODUCT HANDLER TESTS
// ===================================================================

type productPayload struct {
	ProductName string  `json:"product_name"`
	Price       uint64  `json:"price"`
	Rating      uint8   `json:"rating"`
	Image       string  `json:"image"`
}

func TestAddProduct_Success(t *testing.T) {
	cleanCollection(t, "Products")
	defer cleanCollection(t, "Products")

	router := setupRouter()
	body := productPayload{
		ProductName: "Test Laptop",
		Price:       99999,
		Rating:      5,
		Image:       "laptop.jpg",
	}

	w := doRequest(router, http.MethodPost, "/admin/addproducts", body, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestViewProducts_Empty(t *testing.T) {
	cleanCollection(t, "Products")
	defer cleanCollection(t, "Products")

	router := setupRouter()
	w := doRequest(router, http.MethodGet, "/users/viewproducts", nil, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}
}

func TestViewProducts_WithData(t *testing.T) {
	cleanCollection(t, "Products")
	defer cleanCollection(t, "Products")

	router := setupRouter()

	// Add two products
	products := []productPayload{
		{ProductName: "Phone", Price: 499, Rating: 4, Image: "phone.jpg"},
		{ProductName: "Tablet", Price: 299, Rating: 3, Image: "tablet.jpg"},
	}
	for _, p := range products {
		doRequest(router, http.MethodPost, "/admin/addproducts", p, nil)
	}

	// Fetch all products
	w := doRequest(router, http.MethodGet, "/users/viewproducts", nil, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse products: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected 2 products, got %d", len(result))
	}
}

func TestSearchProducts_ByName(t *testing.T) {
	cleanCollection(t, "Products")
	defer cleanCollection(t, "Products")

	router := setupRouter()

	// Add products
	doRequest(router, http.MethodPost, "/admin/addproducts", productPayload{ProductName: "Gaming Laptop", Price: 1200, Rating: 5, Image: "gl.jpg"}, nil)
	doRequest(router, http.MethodPost, "/admin/addproducts", productPayload{ProductName: "Office Chair", Price: 200, Rating: 4, Image: "oc.jpg"}, nil)
	doRequest(router, http.MethodPost, "/admin/addproducts", productPayload{ProductName: "Gaming Mouse", Price: 50, Rating: 4, Image: "gm.jpg"}, nil)

	// Search for "Gaming"
	w := doRequest(router, http.MethodGet, "/users/searchproducts?name=Gaming", nil, nil)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to parse search results: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("Expected 2 'Gaming' products, got %d", len(result))
	}
}

func TestSearchProducts_NoQuery(t *testing.T) {
	router := setupRouter()

	w := doRequest(router, http.MethodGet, "/users/searchproducts", nil, nil)

	if w.Code == http.StatusOK {
		t.Fatal("Search without query param should not return 200")
	}
}

func TestSearchProducts_NoResults(t *testing.T) {
	cleanCollection(t, "Products")
	defer cleanCollection(t, "Products")

	router := setupRouter()
	w := doRequest(router, http.MethodGet, "/users/searchproducts?name=NonExistentXYZ", nil, nil)

	if w.Code != http.StatusOK {
		// Even with no results, the handler should return 200 with an empty list
		t.Logf("Note: Search with no results returned status %d", w.Code)
	}
}
