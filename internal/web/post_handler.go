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

type PostHandler struct {
	store    store.Store
	sessions *scs.SessionManager
}

func (h *PostHandler) createView() http.HandlerFunc {
	type data struct {
		SessionData
		Thread store.Thread
		CSRF   template.HTML // string which is not escaped
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
		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Thread:      t,
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *PostHandler) view() http.HandlerFunc {
	type data struct {
		SessionData
		Thread   store.Thread
		Post     store.Post
		Comments []store.Comment
		CSRF     template.HTML // string which is not escaped
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
		tmpl.Execute(w, data{
			SessionData: GetSessionData(h.sessions, r.Context()),
			Thread:      t,
			Post:        p,
			Comments:    cc,
			CSRF:        csrf.TemplateField(r),
		})
	}
}

func (h *PostHandler) save() http.HandlerFunc {
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
		form := CreatePostForm{
			Title:   r.FormValue("title"),
			Content: r.FormValue("content"),
		}

		if !form.Validate() {
			// lets store the error to the session
			// session cannot store complex types so we need to encode the form see func init in form.go
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}
		//send new post to db
		p := &store.Post{
			ID:       uuid.New(),
			ThreadID: t.ID,
			Title:    form.Title,
			Content:  form.Content,
		}
		if err := h.store.CreatePost(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// add flash message
		h.sessions.Put(r.Context(), "flash", "Your new post has been created.")
		// redirect to the new post
		http.Redirect(w, r, "/threads/"+t.ID.String()+"/"+p.ID.String(), http.StatusFound)
	}
}

func (h *PostHandler) vote() http.HandlerFunc {
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
