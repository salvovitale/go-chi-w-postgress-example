package web

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func NewHandler(s store.Store, ss *scs.SessionManager, csrfKey []byte) *Handler {
	h := &Handler{
		Mux:      chi.NewRouter(),
		store:    s,
		sessions: ss,
	}

	threadsHandler := ThreadHandler{store: s, sessions: ss}
	postHandler := PostHandler{store: s, sessions: ss}
	commentHandler := CommentHandler{store: s, sessions: ss}

	// add logger middleware
	h.Use(middleware.Logger)

	h.Use(csrf.Protect(csrfKey, csrf.Secure(false))) // set security to false for development otherwise the cookie will only be sent over https
	// homepage
	h.Get("/", h.homeView())

	// sub paths
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", threadsHandler.listView())
		r.Get("/new", threadsHandler.createView())
		r.Get("/{id}", threadsHandler.view())
		r.Post("/", threadsHandler.save())
		r.Post("/{id}/delete", threadsHandler.delete())

		// post routes
		r.Get("/{id}/new", postHandler.createView())
		r.Get("/{threadId}/{postId}", postHandler.view())
		r.Get("/{threadId}/{postId}/vote", postHandler.vote())
		r.Post("/{id}", postHandler.save())

		// comment routes
		r.Post("/{threadId}/{postId}", commentHandler.save())
	})

	// comments vote
	h.Get("/comments/{id}/vote", commentHandler.vote())
	return h
}

type Handler struct {
	*chi.Mux //embedded structure
	store    store.Store
	sessions *scs.SessionManager
}

func (h *Handler) homeView() http.HandlerFunc {
	type data struct {
		Posts []store.Post
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve all posts
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Posts: pp})
	}
}
