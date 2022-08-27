package web

import (
	"context"
	"database/sql"
	"encoding/gob"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func init() {
	// we store uuid into the sessions so we need to register the uuid to be encoded and decoded
	gob.Register(uuid.UUID{})
}
func NewSessionManager(dataSourceName string) (*scs.SessionManager, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	sessions := scs.New()
	sessions.Store = postgresstore.New(db)

	return sessions, nil
}

// All sessions data are included in a single struct
type SessionData struct {
	FlashMessage string
	Form         interface{} //so it will work with any form type
	User         store.User
	LoggedIn     bool
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	// Get the flash message from the session, we use pop to remove it from the session because we want to display it only once
	data.FlashMessage = session.PopString(ctx, "flash")
	// retrieve user from context
	data.User, data.LoggedIn = ctx.Value("user").(store.User)

	// Get the form from the session
	data.Form = session.Pop(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{} // initialize the form with an empty map
	}

	return data
}
