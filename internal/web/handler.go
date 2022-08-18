package web

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func NewHandler(s store.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewRouter(),
		store: s,
	}

	// add logger middleware
	h.Use(middleware.Logger)

	// homepage
	h.Get("/", h.homeView())

	// sub paths
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.threadsListView())
		r.Get("/new", h.threadCreateView())
		r.Get("/{id}", h.threadView())
		r.Post("/", h.threadStore())
		r.Post("/{id}/delete", h.threadDelete())

		// post routes
		r.Get("/{id}/new", h.postCreateView())
		r.Get("/{threadId}/{postId}", h.postView())
		r.Post("/{id}", h.postStore())
	})

	return h
}

type Handler struct {
	*chi.Mux
	store store.Store
}

func (h *Handler) homeView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) threadsListView() http.HandlerFunc {
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

func (h *Handler) threadCreateView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) threadView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) threadStore() http.HandlerFunc {
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

func (h *Handler) threadDelete() http.HandlerFunc {
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

// Post routes below here

func (h *Handler) postCreateView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) postView() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) postStore() http.HandlerFunc {
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
