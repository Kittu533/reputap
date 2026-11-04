.PHONY: run build migrate-up migrate-down seed clean

# Go binary path (fallback to 'go' if in PATH)
GO := go

run:
	$(GO) run cmd/server/main.go

build:
	$(GO) build -o bin/server cmd/server/main.go

migrate-up:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-down:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

seed:
	$(GO) run cmd/seed/main.go

clean:
	rm -rf bin/

tidy:
	$(GO) mod tidy
