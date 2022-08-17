.PHONY: postgres adminer migrate migrate-down

postgres:
	docker run --rm -ti -p 5432:5432 -e POSTGRES_PASSWORD=secret postgres

adminer:
	docker run --rm -ti -p 8080:8080 adminer

migrate:
	./ci/db/migrations/migrate -source file://ci/db/migrations \
											 -database postgres://postgres:secret@localhost/postgres?sslmode=disable up

migrate-down:
	./ci/db/migrations/migrate -source file://ci/db/migrations \
											 -database postgres://postgres:secret@localhost/postgres?sslmode=disable down