.PHONY: build run dev clean docker-build docker-run generate generate-lazy migrate-up migrate-down migrate-status migrate-create

build: generate-lazy
	go build -o bin/main cmd/app/main.go

run:
	./bin/main

dev: generate-lazy
	air

clean:
	rm -rf bin/*

# build image (no need for docker-compose since it's just one docker image that Cloud Run will use)
docker-build:
	docker build -t whisp .

# runs the whisp image, mapping port 8080 and mounting the .env file (very basic)
docker-run:
	docker run -p 8080:8080 -v $(PWD)/.env:/app/.env whisp

# tool is either sqlc or templ (too lazy to validate tool arg)
generate:
	$(tool) generate

generate-lazy:
	$(MAKE) generate tool=templ
	$(MAKE) generate tool=sqlc

migrate-up:
	./scripts/migrate.sh up $(db)

migrate-down:
	./scripts/migrate.sh down $(db)

migrate-status:
	./scripts/migrate.sh status $(db)

migrate-create:
	./scripts/migrate.sh create $(db) $(name)
