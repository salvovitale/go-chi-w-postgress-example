.PHONY: postgres adminer migrate migrate-down

postgres:
	docker run --rm -ti -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres

migrate:
	./ci/db/migrations/migrate -source file://ci/db/migrations \
											 -database postgres://postgres:secret@localhost/postgres?sslmode=disable up

migrate-down:
	./ci/db/migrations/migrate -source file://ci/db/migrations \
											 -database postgres://postgres:secret@localhost/postgres?sslmode=disable down

run-w-reflex:
	reflex -s go run ./cmd/server/main.go