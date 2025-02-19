include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
Command := $(firstword $(MAKECMDGOALS))

%::
	@true

live/go:
	@go run github.com/air-verse/air@v1.52.3 \
		--build.cmd "go build -o tmp/main" --build.bin "tmp/main" --build.delay "1000" \
		--build.exclude_dir "node_modules,tmp,vendor" \
		--build.include_ext "go,yaml" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

live/templ:
	@templ generate -watch -proxy="http://127.0.0.1:2468" --open-browser=false

live/tailwind:
	@npx tailwindcss -i input.css -o public/tailwind.css --watch

live: 
	make -j4 live/templ live/tailwind live/go build/commit-id

build/tailwind:
	npx tailwindcss -i input.css -o public/tailwind.css --minify

build/templ:
	templ generate

build/commit-id:
	echo "{\"commit_id\": \"$(Arguments)\"}" > config/commit-id.json

install:
	brew install golang-migrate gnupg
	go install github.com/air-verse/air@v1.52.3
	go install github.com/a-h/templ/cmd/templ@v0.2.793
	npm install
	npx husky init

	@echo 'add "go.goroot:"$$GOROOT" to settings.json VsCode'
	@echo 'after run "make db_init' 



# DB commands
db/encode:
	tar -czvf database/migration/database.tar.gz database/migration
	gpg -c database/migration/database.tar.gz
	shred -u database/migration/*.sql database/migration/database.tar.gz

db/decode:
	gpg database/migration/database.tar.gz.gpg
	tar -xzvf database/migration/database.tar.gz
	shred -u database/migration/database.tar.gz.gpg database/migration/database.tar.gz


db/combine/script:
	cd database/migration/ && cat $$(ls | grep .up.sql)| grep -v '^--' > ../../output.sql

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
migration_create: 
	@migrate create -ext sql -dir database/migration/ -seq $(Arguments)

migration_up: 
	migrate -path database/migration/ -database "${DATABASE_URL}" -verbose up

migration_down: 
	migrate -path database/migration/ -database "${DATABASE_URL}" -verbose down 1

migration_fix: 
	migrate -path database/migration/ -database "${DATABASE_URL}" force ${VERSION}

show_deadcode:
	deadcode .