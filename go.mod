module chirpy

go 1.23.0

toolchain go1.23.6

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	golang.org/x/crypto v0.34.0 // indirect
)

require internal/database v1.0.0

replace internal/database => ./internal/database

require (
	github.com/joho/godotenv v1.5.1
	internal/auth v1.0.0
)

replace internal/auth => ./internal/auth
