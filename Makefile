export PATH := $(HOME)/.bun/bin:$(PATH)

.PHONY: dev dev-be dev-fe db migrate

db:
	docker compose up -d

migrate:
	@for f in $$(ls -1 reputap-be/migrations/*.up.sql | sort); do \
		docker exec -i reputap-mysql mysql -ureputap -psecret reputap < "$$f"; \
	done

dev-be:
	cd reputap-be && go run cmd/server/main.go

dev-fe:
	cd reputap-fe && bun run dev

dev:
	@$(MAKE) -j 2 dev-be dev-fe
