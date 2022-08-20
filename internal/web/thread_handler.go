package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

type ThreadHandler struct {
	store store.Store
}

func (h *ThreadHandler) listView() http.HandlerFunc {
	// wrap some local data that wont be visible from outside
	type data struct {
		Threads []store.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Threads: threads})
	}
}

func (h *ThreadHandler) createView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *ThreadHandler) view() http.HandlerFunc {
	type data struct {
		Thread store.Thread
		Posts  []store.Post
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
		tmpl.Execute(w, data{Thread: t, Posts: pp})
	}
}

func (h *ThreadHandler) save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the form
		title := r.FormValue("title")
		description := r.FormValue("description")
		//TODO validate the form
		//send new thread to db
		if err := h.store.CreateThread(&store.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
		// redirect to the thread list
		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
