package web

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

type CommentHandler struct {
	store    store.Store
	sessions *scs.SessionManager
}

func (h *CommentHandler) save() http.HandlerFunc {
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
		form := CreateCommentForm{
			Content: r.FormValue("content"),
		}

		if !form.Validate() {
			h.sessions.Put(r.Context(), "form", form)
			http.Redirect(w, r, r.Referer(), http.StatusFound)
			return
		}

		//send new comment to db
		if err := h.store.CreateComment(&store.Comment{
			ID:      uuid.New(),
			PostID:  p.ID,
			Content: form.Content,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// add flash message
		h.sessions.Put(r.Context(), "flash", "Your comment has been submitted.")

		// redirect to the new post
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}
}

func (h *CommentHandler) vote() http.HandlerFunc {
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
