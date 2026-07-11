# 🛒 Ecommerce Backend API

A comprehensive, production-ready RESTful API for an e-commerce platform built with **Go**, the **Gin** web framework, and **MongoDB**. This project showcases modern backend engineering practices including **Docker containerization**, a full suite of **Unit & Integration tests**, **Swagger (OpenAPI) documentation**, **Role-Based Authentication**, and **API Request Logging & Monitoring**.

---

## ✨ Features

| Category | Capabilities |
|---|---|
| **Auth** | User signup & login with JWT (access + refresh tokens), bcrypt password hashing |
| **Products** | Admin product creation, list all products, regex-powered search |
| **Cart** | Add / remove items, view cart with aggregated total price |
| **Orders** | Checkout entire cart or instant-buy a single product |
| **Addresses** | Add up to 2 addresses per user (home / work), edit & delete |
| **Documentation** | Interactive OpenAPI 3.0 specs available via **Swagger UI** |
| **Testing** | 30+ comprehensive unit and integration tests using a dedicated test database |
| **DevOps** | Multi-stage Docker builds and Docker Compose orchestration |
| **Monitoring** | Built-in HTTP request logging (status codes, latency, client IP) via Gin middleware |
---

## 🏗️ Tech Stack

- **Language** — Go 1.24
- **Framework** — [Gin](https://github.com/gin-gonic/gin)
- **Database** — MongoDB (Atlas or self-hosted)
- **Auth** — JWT via `dgrijalva/jwt-go`
- **Validation** — `go-playground/validator/v10`
- **Password Hashing** — `golang.org/x/crypto/bcrypt`
- **Config** — `.env` loaded with `godotenv`

---

## 📁 Project Structure

```
Backend/
├── main.go                          # Application entry-point & route registration
├── controller/
│   ├── controller.go                # Auth (signup/login) & product handlers
│   ├── cart.go                      # Cart & order handlers
│   ├── address.go                   # Address CRUD handlers
│   ├── testhelpers_test.go          # Test infrastructure (TestMain, helpers)
│   └── controller_test.go           # Unit + integration tests for handlers
├── database/
│   ├── databasesetup.go             # MongoDB connection & collection helpers
│   └── cart.go                      # Cart / order database operations
├── middleware/
│   └── middleware.go                # JWT authentication middleware
├── models/
│   └── models.go                    # User, Product, Order, Address, Payment models
├── routes/
│   └── routes.go                    # Public route definitions
├── tokens/
│   ├── tokengen.go                  # JWT token generation, validation & refresh
│   └── tokengen_test.go             # Token unit tests
├── Dockerfile                       # Multi-stage Docker build
├── docker-compose.yaml              # Docker Compose (MongoDB + Backend)
├── .dockerignore                    # Docker build context exclusions
├── .env                             # Environment variables (git-ignored)
├── go.mod / go.sum                  # Go module dependencies
└── readme.md                        # ← You are here
```

---

## 🚀 Getting Started

### Prerequisites

| Requirement | Version | Required? |
|---|---|---|
| **Go** | ≥ 1.24 | ✅ For local development |
| **MongoDB** | ≥ 6.0 | ✅ Local install **or** via Docker |
| **Docker** & **Docker Compose** | Latest | ⚡ Optional — for containerised run |

### 1. Clone the repository

```bash
git clone https://github.com/Hifzu04/Ecommerce-Backend-with-Golang.git
cd Ecommerce-Backend-with-Golang/Backend
```

### 2. Configure environment

Create a `.env` file in the `Backend/` directory:

```env
PORT=8000
MONGO_URI=mongodb://localhost:27017
SECRET_LOVE=my-super-secret-key-change-in-production
```

> **Note:** If using MongoDB Atlas instead of local, set `MONGO_URI=mongodb+srv://<user>:<pass>@<cluster>.mongodb.net/`

### 3a. Run locally (without Docker)

Make sure MongoDB is running on your machine, then:

```bash
# Install dependencies
go mod download

# Start the server
go run main.go
```

You should see:
```
Connected to mongodb!!
[GIN-debug] Listening and serving HTTP on :8000
```

The server is now live at **`http://localhost:8000`**.

### 📖 Accessing Swagger API Documentation

This project uses **Swagger** to provide interactive OpenAPI documentation.

Once the server is running (either locally or via Docker), open your browser and navigate to:
**👉 `http://localhost:8000/swagger/index.html`**

You can use the Swagger UI to view all endpoints, expected payloads, and even execute real HTTP requests against the local API.

**Verify it works:**
```bash
curl http://localhost:8000/users/viewproducts
```

### 3b. Run with Docker Compose

Docker Compose spins up both **MongoDB** and the **Go backend** automatically — no local MongoDB needed:

```bash
# Build & start both containers
docker compose up --build -d

# View live logs
docker compose logs -f backend

# Stop everything
docker compose down
```

The API is available at **`http://localhost:8000`**.

> **Tip:** If port 27017 is already in use (local MongoDB running), either stop it with `sudo systemctl stop mongod` or change the port mapping in `docker-compose.yaml`.

---

## 🐳 Docker

### Architecture

```
┌─────────────────────────────────────────┐
│           docker-compose.yaml           │
│                                         │
│  ┌─────────────┐    ┌────────────────┐  │
│  │  MongoDB 7.0 │◄───│  Go Backend    │  │
│  │  :27017      │    │  :8000         │  │
│  └─────────────┘    └────────────────┘  │
│         │                    │           │
│    mongo-data           Dockerfile       │
│    (volume)          (multi-stage)       │
└─────────────────────────────────────────┘
```

### Dockerfile — Multi-stage build

| Stage | Base Image | Purpose |
|---|---|---|
| **builder** | `golang:1.24-alpine` | Downloads deps & compiles a static binary |
| **runtime** | `alpine:3.20` | Runs the binary (~15 MB final image) |

### Build manually (without Compose)

```bash
docker build -t ecommerce-backend .
docker run -p 8000:8000 --env-file .env ecommerce-backend
```

---

## 🧪 Testing

The project includes **30 automated tests** — unit tests for core logic and integration tests for HTTP handlers.

### Running Tests

> **Prerequisite:** MongoDB must be running on `localhost:27017`. Tests use a separate `TestEcommerce` database that is automatically cleaned up after each run — your production data is never touched.

```bash
cd Backend/

# Run ALL tests
go test ./... -v

# Run only controller tests (password + handler tests)
go test ./controller/... -v

# Run only token tests (JWT generation & validation)
go test ./tokens/... -v

# Run a single specific test
go test ./controller/... -v -run TestSignupHandler_Success

# Run with coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out    # Opens report in browser
```

### Test Summary

| Package | Type | Tests | What's Covered |
|---|---|---|---|
| `controller` | Unit | 6 | Password hashing, bcrypt salt randomness, verification |
| `controller` | Integration | 14 | Signup (success, duplicate email/phone, invalid body), Login (success, wrong password, non-existent user), Products (add, view, search) |
| `tokens` | Unit | 10 | Token generation, validation, tampering detection, expiry, wrong secret, claims integrity |

### Expected Output

```
=== RUN   TestHashPassword
--- PASS: TestHashPassword (1.84s)
=== RUN   TestHashPassword_DifferentInputs_DifferentHashes
--- PASS: TestHashPassword_DifferentInputs_DifferentHashes (3.66s)
...
=== RUN   TestSignupHandler_Success
--- PASS: TestSignupHandler_Success (1.83s)
=== RUN   TestSignupHandler_DuplicateEmail
--- PASS: TestSignupHandler_DuplicateEmail (1.83s)
...
=== RUN   TestTokenGenerator_Success
--- PASS: TestTokenGenerator_Success (0.00s)
=== RUN   TestValidateToken_ExpiredToken
--- PASS: TestValidateToken_ExpiredToken (0.00s)
...
ok   github.com/Hifzu04/Ecommerce/Backend/controller  32.640s
ok   github.com/Hifzu04/Ecommerce/Backend/tokens       0.014s
```

---

## 📡 API Endpoints

> Routes **above** the authentication middleware are **public**; routes **below** it require a `token` header with a valid JWT.

### Public Routes

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/users/signup` | Register a new user |
| `POST` | `/users/login` | Login & receive JWT tokens |
| `POST` | `/admin/addproducts` | Add a new product (admin) |
| `GET` | `/users/viewproducts` | List all products |
| `GET` | `/users/searchproducts?name=<query>` | Search products by name (regex) |

### Protected Routes *(require `token` header)*

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/addtocart?id=<productID>&userID=<userID>` | Add product to cart |
| `GET` | `/removeitem?id=<productID>&userID=<userID>` | Remove product from cart |
| `GET` | `/listcart?id=<userID>` | View cart items & total |
| `POST` | `/cartcheckout?id=<userID>` | Purchase all cart items |
| `GET` | `/instantbuy?Userid=<userID>&pid=<productID>` | Instant buy a single product |
| `POST` | `/addaddress?id=<userID>` | Add a delivery address |
| `PUT` | `/edithomeaddress?id=<userID>` | Edit home address |
| `PUT` | `/editworkaddress?id=<userID>` | Edit work address |
| `DELETE` | `/deleteaddresses?id=<userID>` | Delete all addresses |

---

## 📋 Sample Requests

### Signup

```bash
curl -X POST http://localhost:8000/users/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Hifzur",
    "last_name": "Rahman",
    "email": "hifzur@example.com",
    "password": "SecurePass123",
    "phone": "+4534545435"
  }'
```

### Login

```bash
curl -X POST http://localhost:8000/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "hifzur@example.com",
    "password": "SecurePass123"
  }'
```

### Add Product (Admin)

```bash
curl -X POST http://localhost:8000/admin/addproducts \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Laptop",
    "price": 99999,
    "rating": 5,
    "image": "laptop.jpg"
  }'
```

### Add to Cart *(protected)*

```bash
curl http://localhost:8000/addtocart?id=<productID>&userID=<userID> \
  -H "token: <your_jwt_token>"
```

### Add Address *(protected)*

```bash
curl -X POST "http://localhost:8000/addaddress?id=<userID>" \
  -H "Content-Type: application/json" \
  -H "token: <your_jwt_token>" \
  -d '{
    "house_name": "123",
    "street_name": "Main St",
    "city_name": "Copenhagen",
    "pin_code": "2100"
  }'
```

---

## 🗄️ Data Models

### User
```
ID, First_Name, Last_Name, Email, Password, Phone,
Token, Refresh_Token, Created_At, Updated_At, User_ID,
UserCart []ProductUser, Address_Details []Address, Order_Status []Order
```

### Product / ProductUser
```
Product_ID, Product_Name, Price, Rating, Image
```

### Order
```
Order_ID, Order_Cart []ProductUser, Ordered_At, Price, Discount, Payment_Method
```

### Address
```
Address_id, House, Street, City, Pincode, State
```

---

## 🔒 Authentication Flow

1. **Signup** → password is hashed with bcrypt → user stored in MongoDB → JWT + refresh token returned.
2. **Login** → email looked up → password verified against hash → new JWT + refresh token generated.
3. **Protected routes** → `token` header validated by middleware → `email` and `uid` injected into Gin context.

---

## ❓ Troubleshooting

| Problem | Solution |
|---|---|
| `address already in use :8000` | Another process is using port 8000. Kill it: `lsof -i :8000` then `kill <PID>`, or change `PORT` in `.env` |
| `address already in use :27017` (Docker) | Local MongoDB is running. Stop it: `sudo systemctl stop mongod` |
| `error parsing uri: no such host` | Your `MONGO_URI` in `.env` is pointing to a dead/unreachable cluster. Use `mongodb://localhost:27017` for local MongoDB |
| `MONGO_URI environment variable is not set` | Create a `.env` file in `Backend/` with the required variables (see step 2 above) |
| Tests fail with `no such host` | MongoDB isn't running. Start it: `sudo systemctl start mongod` or `docker compose up -d mongodb` |

---

## 📄 License

This project is open-source and available under the [MIT License](LICENSE).
