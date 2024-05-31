include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

tailwind:
	./tailwindcss -i input.css -o public/tailwind.css --watch

templ:
	templ generate -watch -proxy=http://localhost:2468

install:
	brew install golang-migrate

	# GO PACKAGES
	go get -u github.com/go-chi/chi/v5
	go get golang.org/x/crypto/bcrypt
	go get github.com/joho/godotenv
	go get github.com/go-sql-driver/mysql
	go get github.com/go-chi/jwtauth/v5

	go get github.com/golang-jwt/jwt/v5

# DB commands

db_init:
	docker-compose up -d

db_start: db_init
	docker-compose start
	
db_stop: db_start
	docker-compose down

# Migration commands

migration_up: 
	migrate -path database/migration/ -database "${DATABASE_URL}" -verbose up

migration_down: 
	migrate -path database/migration/ -database "${DATABASE_URL}" -verbose down 1

migration_fix: 
	migrate -path database/migration/ -database "${DATABASE_URL}" force ${VERSION}
