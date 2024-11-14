include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

tailwind:
	./tailwindcss -i input.css -o public/tailwind.css --watch

templ:
	templ generate -watch

air:
	air -c .air.toml

install:
	brew install golang-migrate
	go install github.com/air-verse/air@v1.52.3
	go install github.com/a-h/templ/cmd/templ@v0.2.778

	@echo 'add "alias air=$$GOPATH/bin/air" in .bashrc / .zshrc' 

start: 
	@@ ./tailwindcss -i input.css -o public/tailwind.css --watch & \
	templ generate -watch & \
	air -c .air.toml

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

show_deadcode:
	deadcode .