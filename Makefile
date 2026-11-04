export PATH := $(HOME)/.bun/bin:$(PATH)

.PHONY: dev dev-be dev-fe db migrate seed produce build-linux

db:
	docker compose up -d postgres

migrate:
	cd reputap-be && go run cmd/migrate/main.go

build-linux:
	cd reputap-be && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/server ./cmd/server
	cd reputap-be && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/migrate ./cmd/migrate

seed:
	cd reputap-be && go run cmd/seed/main.go

produce:
	cd reputap-be && go run cmd/produce/main.go -n 20

dev-be:
	cd reputap-be && go run cmd/server/main.go

dev-fe:
	cd reputap-fe && bun run dev

dev:
	@$(MAKE) -j 2 dev-be dev-fe
