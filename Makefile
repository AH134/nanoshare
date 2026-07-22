include .env
export
.PHONY: dev-backend dev-frontend clean migrate-up migrate-down create-admin reset-db

dev-backend:
	air

dev-frontend:
	npm --prefix ui run dev

migrate-up:
	goose -dir ./internal/database/migrations sqlite3 ./data/nanoshare.db up

migrate-down:
	goose -dir ./internal/database/migrations sqlite3 ./data/nanoshare.db down

create-admin:
	go run ./cmd/nanoshare/createadmin/main.go --username=$(ADMIN_USERNAME) --password=$(ADMIN_PASSWORD)

clean:
	rm -rf air/tmp ui/dist data/*
