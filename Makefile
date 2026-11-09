export PATH := $(HOME)/.bun/bin:$(PATH)

.PHONY: dev dev-be dev-fe db migrate seed produce reset-db build-linux

db:
	docker compose up -d postgres

migrate:
	cd reputap-be && go run cmd/migrate/main.go

reset-db:
	cd reputap-be && go run cmd/reset/main.go

build-linux:
	cd reputap-be && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/server ./cmd/server
	cd reputap-be && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/migrate ./cmd/migrate
	cd reputap-be && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ../bin/produce ./cmd/produce

seed:
	cd reputap-be && go run cmd/seed/main.go

produce:
	cd reputap-be && go run cmd/produce/main.go -n 20 -prefix RP26 -prod

produce-local:
	cd reputap-be && go run cmd/produce/main.go -n 20 -prefix RP26

dev-be:
	cd reputap-be && go run cmd/server/main.go

dev-fe:
	cd reputap-fe && bun run dev

dev:
	@$(MAKE) -j 2 dev-be dev-fe
