package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func NewPostStore(db *sqlx.DB) *PostStore {
	return &PostStore{DB: db}
}

type PostStore struct {
	*sqlx.DB
}

func (s *PostStore) PostsByThread(threadID uuid.UUID) ([]store.Post, error) {
	var p []store.Post
	var query = `
		SELECT
			posts.*,
			COUNT(comments.*) AS comments_count
		FROM posts
		JOIN comments ON comments.post_id = posts.id
		WHERE thread_id = $1
		GROUP BY posts.id
		ORDER BY votes DESC
	`
	if err := s.Select(&p, query, threadID); err != nil {
		return []store.Post{}, fmt.Errorf("error getting posts: %w", err)
	}
	return p, nil
}

func (s *PostStore) Posts() ([]store.Post, error) {
	var p []store.Post
	var query = `
		SELECT
			posts.*,
			COUNT(comments.*) AS comments_count,
			threads.title AS thread_title
		FROM posts
		JOIN comments ON comments.post_id = posts.id
		JOIN threads ON threads.id = posts.thread_id
		GROUP BY posts.id, threads.title
		ORDER BY votes DESC
	`
	if err := s.Select(&p, query); err != nil {
		return []store.Post{}, fmt.Errorf("error getting posts: %w", err)
	}
	return p, nil
}

func (s *PostStore) Post(id uuid.UUID) (store.Post, error) {
	var p store.Post
	if err := s.Get(&p, "SELECT * FROM posts WHERE id = $1", id); err != nil {
		return store.Post{}, fmt.Errorf("error getting post: %w", err)
	}
	return p, nil
}

func (s *PostStore) CreatePost(p *store.Post) error {
	if err := s.Get(p, "INSERT INTO posts VALUES ($1, $2, $3, $4, $5) RETURNING *",
		p.ID,
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes); err != nil {
		return fmt.Errorf("error creating post: %w", err)
	}
	return nil
}

func (s *PostStore) UpdatePost(p *store.Post) error {
	if err := s.Get(p, "UPDATE posts SET thread_id = $1, title = $2, content = $3, votes = $4 WHERE id = $5 RETURNING *",
		p.ThreadID,
		p.Title,
		p.Content,
		p.Votes,
		p.ID); err != nil {
		return fmt.Errorf("error updating post: %w", err)
	}
	return nil
}

func (s *PostStore) DeletePost(id uuid.UUID) error {
	if _, err := s.Exec("DELETE FROM posts WHERE id = $1", id); err != nil {
		return fmt.Errorf("error deleting post: %w", err)
	}
	return nil
}
