package web

import (
	"context"
	"html/template"
	"net/http"
	"sync"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
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
	userHandler := UserHandler{store: s, sessions: ss}

	// add logger middleware
	h.Use(middleware.Logger)

	// add csrf protection middleware
	h.Use(csrf.Protect(csrfKey, csrf.Secure(false))) // set security to false for development otherwise the cookie will only be sent over https

	// add session middleware
	h.Use(ss.LoadAndSave)

	// add custom middleware to retrieve the user from the session and add it to the request context
	h.Use(h.withUser)

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

	// user routes
	h.Get("/register", userHandler.RegisterView())
	h.Post("/register", userHandler.Register())
	h.Get("/login", userHandler.LoginView())
	h.Post("/login", userHandler.Login())
	h.Get("/logout", userHandler.Logout())
	return h
}

type Handler struct {
	*chi.Mux //embedded structure
	store    store.Store
	sessions *scs.SessionManager
}

func (h *Handler) homeView() http.HandlerFunc {
	type data struct {
		SessionData
		Posts []store.Post
	}

	var once sync.Once
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		// retrieve all posts
		pp, err := h.store.Posts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		once.Do(func() {
			h.sessions.Put(r.Context(), "flash", "Welcome!")
		})

		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Posts:       pp,
		})
	}
}

// create a middleware to retrieve the user from the session and add it to the request context
func (h *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, _ := h.sessions.Get(r.Context(), "user_id").(uuid.UUID)

		user, err := h.store.User(id)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
