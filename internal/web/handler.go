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
	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsListView())
		r.Get("/new", h.ThreadsCreateView())
		r.Post("/", h.ThreadsStore())
		r.Post("/{id}/delete", h.ThreadsDelete())
	})

	h.Get("/html", func(w http.ResponseWriter, r *http.Request) {
		// "New" should contain the name of the root where the includes will be added
		t := template.Must(template.New("layout.html").ParseGlob("templates/includes/*.html"))
		// now we can add the rest of the file in "t"
		t = template.Must(t.ParseFiles("templates/layout.html", "templates/child-template.html"))
		type params struct {
			Title   string
			Text    string
			Lines   []string
			Number1 int
			Number2 int
		}

		t.Execute(w, params{
			Title: "Reddit Clone",
			Text:  "Welcome to the Reddit Clone",
			Lines: []string{
				"Lines1",
				"Lines2",
				"Lines3",
			},
			Number1: 1,
			Number2: 2,
		})
	})

	return h
}

type Handler struct {
	*chi.Mux
	store store.Store
}

const threadsListHTML = `
	<h1>Threads 2</h1>
	<dl>
	{{range .Threads}}
		<dt><strong>{{.Title}}</strong></dt>
		<dd>{{.Description}}</dd>
		<dd>
			<form action="/threads/{{.ID}}/delete" method="POST">
				<button type="submit">Delete</button>
			</form>
		</dd>
	{{end}}
	</dl>
	<a href="/threads/new">Create Thread</a>
`

func (h *Handler) ThreadsListView() http.HandlerFunc {
	// wrap some local data that wont be visible from outside
	type data struct {
		Threads []store.Thread
	}

	// this will be run once when the server starts
	// if the parsing fails the application wont boot as it would be unable to render the page
	tmpl := template.Must(template.New("").Parse(threadsListHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		threads, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Threads: threads})
	}
}

const threadCreateHTML = `
	<h1>New Thread</h1>
	<form action="/threads" method="post">
		<table>
			<tr>
				<td>Title</td>
				<td><input type="text" name="title" /></td>
			</tr>
			<tr>
				<td>Description</td>
				<td><input type="text" name="description" /></td>
			</tr>
		</table>
		<button type="submit">Create thread</button>
	</form>
`

func (h *Handler) ThreadsCreateView() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(threadCreateHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) ThreadsStore() http.HandlerFunc {
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

func (h *Handler) ThreadsDelete() http.HandlerFunc {
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
