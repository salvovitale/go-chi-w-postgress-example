package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/csrf"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

type ThreadHandler struct {
	store    store.Store
	sessions *scs.SessionManager
}

func (h *ThreadHandler) listView() http.HandlerFunc {
	// wrap some local data that wont be visible from outside
	type data struct {
		SessionData
		Threads []store.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Threads:     threads,
		})
	}
}

func (h *ThreadHandler) createView() http.HandlerFunc {
	type data struct {
		SessionData
		CSRF template.HTML // string which is not escaped
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *ThreadHandler) view() http.HandlerFunc {
	type data struct {
		SessionData
		Thread store.Thread
		Posts  []store.Post
		CSRF   template.HTML // string which is not escaped
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id
		idStr := chi.URLParam(r, "id")
		//parse and validate the id
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pp, err := h.store.PostsByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Thread:      t,
			Posts:       pp,
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *ThreadHandler) save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the form
		form := CreateThreadForm{
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
		}

		if !form.Validate() {
			// lets store the error to the session
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		//send new thread to db
		if err := h.store.CreateThread(&store.Thread{
			ID:          uuid.New(),
			Title:       form.Title,
			Description: form.Description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// add flash message
		h.sessions.Put(r.Context(), "flash", "Your new thread has been created.")
		// redirect to the thread list
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}

func (h *ThreadHandler) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the id
		idStr := chi.URLParam(r, "id")
		//parse and validate the id
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//delete thread from db
		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// add flash message
		h.sessions.Put(r.Context(), "flash", "Your thread has been deleted.")

		// redirect to the thread list
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
