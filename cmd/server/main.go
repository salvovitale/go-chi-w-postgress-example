package main

import (
	"log"
	"net/http"

	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/postgres"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/web"
)

func main() {

	dsn := "postgres://postgres:secret@localhost/postgres?sslmode=disable"

	store, err := postgres.NewStore(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sessions, err := web.NewSessionManager(dsn)
	if err != nil {
		log.Fatal(err)
	}

	csrfKey := []byte("01234567890123456789012345678901") //32 bytes long
	h := web.NewHandler(store, sessions, csrfKey)

	// to avoid the error scs: no session data in context we need to wrap the web handler which in this case embeds the chi mux into the LoadAndSave middleware
	http.ListenAndServe(":3000", sessions.LoadAndSave(h))
}
