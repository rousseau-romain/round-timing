include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

db:
	echo ${DATABASE_URL}

db/combine/script:
	cd database/migration/ && cat $(ls | grep .up.sql) > output.sql

live/templ:
	templ generate --watch --proxy="http://127.0.0.1:2468" --cmd="go run ."

live/tailwind:
	@npx tailwindcss -i input.css -o public/tailwind.css --watch

live: 
	make -j2 live/templ live/tailwind

build/tailwind:
	npx tailwindcss -i input.css -o public/tailwind.css --minify

build/templ:
	templ generate

install:
	brew install golang-migrate
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	npm install

	@echo 'add "go.goroot:"$$GOROOT" to settings.json VsCode'
	@echo 'after run "make live' 



# DB commands
db_init:
	docker-compose up -d
	echo "Wait 2s"
	@sleep 3
	make migration_up

db_start:
	docker-compose start
	
db_stop:
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