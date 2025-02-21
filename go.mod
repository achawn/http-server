module chirpy

go 1.22.4

require (
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)

require (
	github.com/joho/godotenv v1.5.1 // indirect
	golang.org/x/crypto v0.33.0 // indirect
)

require internal/database v1.0.0

replace internal/database => ./internal/database

require internal/auth v1.0.0

replace internal/auth => ./internal/auth
