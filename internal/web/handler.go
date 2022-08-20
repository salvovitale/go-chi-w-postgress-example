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
		r.Get("/{threadId}/{postId}/vote", h.postVote())
		r.Post("/{id}", h.postStore())

		// comment routes
		r.Post("/{threadId}/{postId}", h.commentStore())
	})

	// comments vote
	h.Get("/comments/{id}/vote", h.commentVote())
	return h
}

type Handler struct {
	*chi.Mux
	store store.Store
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
	type data struct {
		Thread store.Thread
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
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
		tmpl.Execute(w, data{Thread: t})
	}
}

func (h *Handler) postView() http.HandlerFunc {
	type data struct {
		Thread   store.Thread
		Post     store.Post
		Comments []store.Comment
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the thread id to which the Post belongs
		threadIdStr := chi.URLParam(r, "threadId")

		//parse and validate the id
		threadId, err := uuid.Parse(threadIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// retrieve the thread from db
		t, err := h.store.Thread(threadId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//parse post id
		postIdStr := chi.URLParam(r, "postId")

		//parse and validate the id
		postId, err := uuid.Parse(postIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// retrieve the post from db
		p, err := h.store.Post(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// retrieve the comments from db
		cc, err := h.store.CommentsByPost(p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// execute the template passing both the thread and post
		tmpl.Execute(w, data{Thread: t, Post: p, Comments: cc})
	}
}

func (h *Handler) postStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the thread id to which the new Post will belong
		idStr := chi.URLParam(r, "id")

		//parse and validate the id
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// verify that the thread exists
		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//parse the form for new post info
		title := r.FormValue("title")
		content := r.FormValue("content")
		//TODO validate the form
		//send new post to db
		p := &store.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    title,
			Content:  content,
		}
		if err := h.store.CreatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// redirect to the new post
		http.Redirect(w, r, "/threads/"+t.ID.String()+"/"+p.ID.String(), http.StatusFound)
	}
}

func (h *Handler) postVote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the post id
		idStr := chi.URLParam(r, "postId")

		//parse and validate the id
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// retrieve the post
		p, err := h.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//increase/decrease the vote
		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			p.Votes++
		} else if dir == "down" {
			p.Votes--
		}

		//update the comment in db
		if err := h.store.UpdatePost(&p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// redirect to the same page
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *Handler) commentStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse post id
		postIdStr := chi.URLParam(r, "postId")

		//parse and validate the id
		postId, err := uuid.Parse(postIdStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// retrieve the post from db to verify that it exists
		p, err := h.store.Post(postId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//parse the form for new comment info
		content := r.FormValue("content")
		//TODO validate the form
		//send new comment to db
		if err := h.store.CreateComment(&store.Comment{
			ID:      uuid.New(),
			PostID:  p.ID,
			Content: content,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// redirect to the new post
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *Handler) commentVote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//parse the comment id
		idStr := chi.URLParam(r, "id")

		//parse and validate the id
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// retrieve the comment
		c, err := h.store.Comment(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//increase/decrease the vote
		dir := r.URL.Query().Get("dir")
		if dir == "up" {
			c.Votes++
		} else if dir == "down" {
			c.Votes--
		}

		//update the comment in db
		if err := h.store.UpdateComment(&c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// redirect to the same page
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}
