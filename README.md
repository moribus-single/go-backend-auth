# Auth service (part)
Part of auth service for generating and refreshing access (JWT) and refresh (RT) token pair. <br> <br>
JWT config: algorithm SHA512 <br>
RT config: algorithm SHA256 (valid base64 string)

## Routes
- **GET /token** - get token pair (JWT, RT) by guid query parameter (requires data in MongoDB)
- **POST /refresh** - refresh token pair (JWT, RT) by refresh token and guid (also requires data in MongoDB)

## Packages
- go.mongodb.org/mongo-driver/mongo
- golang.org/x/crypto/bcrypt
- github.com/golang-jwt/jwt
- github.com/joho/godotenv

## Usage
1. Run MongoDB, create database and table, add some data into table.
2. Set env variables according `.env.example` in `.env` file.
3. Install required packages
4. Run the app.
