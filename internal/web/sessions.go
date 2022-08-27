package web

import (
	"context"
	"database/sql"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
)

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
	// UserId uuid.UUID
}

func GetSessionData(session *scs.SessionManager, ctx context.Context) SessionData {
	var data SessionData

	// Get the flash message from the session, we use pop to remove it from the session because we want to display it only once
	data.FlashMessage = session.PopString(ctx, "flash")
	// data.UserId, _ = session.Get(ctx, "user_id").(uuid.UUID)

	// Get the form from the session
	data.Form = session.Pop(ctx, "form")
	if data.Form == nil {
		data.Form = map[string]string{} // initialize the form with an empty map
	}

	return data
}
