include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

live/templ:
	@templ generate --watch --proxy="http://localhost:2468" --open-browser=false -v

live/tailwind:
	@./tailwindcss -i input.css -o public/tailwind.css --minify --watch

live/air:
	@air -c .air.toml

live: 
	@make -j3 live/templ live/air live/tailwind 

install:
	brew install golang-migrate
	go install github.com/air-verse/air@v1.52.3
	go install github.com/a-h/templ/cmd/templ@v0.2.793

	@echo 'add "alias air=$$GOPATH/bin/air" in .bashrc / .zshrc' 
	@echo 'after run "make live' 

# DB commands
db_init:
	docker-compose up -d
	make migration_up

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