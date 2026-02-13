include .env

DATABASE_URL=${DB_DRIVER}://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}

Arguments := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
Command := $(firstword $(MAKECMDGOALS))

%::
	@true

live/go:
	@go run github.com/air-verse/air@v1.64.0 \
		--build.cmd "go build -o tmp/main" --build.entrypoint "tmp/main" --build.delay "1000" \
		--build.exclude_dir "node_modules,tmp,vendor" \
		--build.include_ext "go,yaml" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

live/templ:
	@go tool templ generate -watch -proxy="http://127.0.0.1:2468" -cmd="./scripts/run-with-log.sh" --open-browser=false

live/tailwind:
	@npx tailwindcss -i input.css -o public/tailwind.css --watch

live:
	make -j3 live/templ live/tailwind

lint/tailwind:
	npx rustywind --check-formatted views/

fix/tailwind:
	npx rustywind --write views/

build/tailwind:
	npx tailwindcss -i input.css -o public/tailwind.css --minify

build/templ:
	go tool templ generate

install:
	go mod tidy
	npm install

# DB commands
db/encrypt:
	tar -czvf database/migration/database.tar.gz database/migration
	gpg -c database/migration/database.tar.gz
	shred -u database/migration/*.sql database/migration/database.tar.gz

db/decrypt:
	gpg database/migration/database.tar.gz.gpg
	tar -xzvf database/migration/database.tar.gz
	shred -u database/migration/database.tar.gz.gpg database/migration/database.tar.gz

db/combine/script:
	cd database/migration/ && cat $$(ls | grep .up.sql)| grep -v '^--' | grep -v '^START TRANSACTION;' | grep -v '^COMMIT;' > ../../output.sql

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

# Monitoring commands
monitoring_start:
	docker-compose up -d loki promtail grafana

monitoring_stop:
	docker-compose stop loki promtail grafana