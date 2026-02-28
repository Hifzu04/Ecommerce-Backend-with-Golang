
# Ecommerce Backend API

A Go-based REST API for an e-commerce platform built with Gin framework and MongoDB.

## Features

- **User Management**: Signup, login with JWT authentication
- **Product Management**: Add products (admin), view all products, search by query
- **Shopping Cart**: Add/remove items, view cart, instant buy
- **Orders**: Place orders from cart or instant buy
- **Addresses**: Add, edit, and delete delivery addresses (max 2 per user)

## Tech Stack

- **Language**: Go
- **Web Framework**: Gin
- **Database**: MongoDB
- **Authentication**: JWT tokens
- **Password Security**: bcrypt hashing

## API Endpoints

### Authentication
- `POST /users/signup` - Register new user
- `POST /users/login` - Login user

### Products
- `GET /users/viewproducts` - Get all products
- `GET /users/searchproducts?name=<query>` - Search products
- `POST /admin/addproducts` - Add product (admin only)

### Cart
- `GET /addtocart?id=<productID>&userID=<userID>` - Add to cart
- `GET /removefromcart?id=<productID>&userID=<userID>` - Remove from cart
- `GET /getcart?id=<userID>` - View cart items and total

### Orders
- `GET /buy?id=<userID>` - Purchase from cart
- `GET /instantbuy?Userid=<userID>&pid=<productID>` - Instant purchase

### Addresses
- `POST /addaddress?id=<userID>` - Add delivery address
- `PUT /editaddress/home?id=<userID>` - Edit home address
- `PUT /editaddress/work?id=<userID>` - Edit work address
- `DELETE /deleteaddress?id=<userID>` - Delete all addresses

## Sample Requests

**Signup Example:**
```json
{
    "first_name": "Hifzur",
    "last_name": "Rahman",
    "email": "hifzurRahman@gmail.com",
    "password": "Hifzur",
    "phone": "+4534545435"
}
```

**Login Example:**
```json
{
    "email": "hifzurRahman@gmail.com",
    "password": "Hifzur"
}
```

**Add Product Example:**
```json
{
    "name": "Laptop",
    "price": 999.99,
    "description": "High-performance laptop",
    "image": "laptop.jpg"
}
```

**Add Address Example:**
```json
{
    "address": "123 Main St",
    "city": "Copenhagen",
    "state": "Zip 2100",
    "type": "home"
}
```

