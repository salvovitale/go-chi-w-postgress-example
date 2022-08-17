package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func NewThreadStore(db *sqlx.DB) *ThreadStore {
	return &ThreadStore{DB: db}
}

type ThreadStore struct {
	// embedded structure so we inherit all the methods from it
	*sqlx.DB
}

func (s *ThreadStore) Threads() ([]store.Thread, error) {
	var t []store.Thread
	if err := s.Select(&t, "SELECT * FROM threads"); err != nil {
		return []store.Thread{}, fmt.Errorf("error getting threads: %w", err)
	}
	return t, nil
}

func (s *ThreadStore) Thread(id uuid.UUID) (store.Thread, error) {
	var t store.Thread
	if err := s.Get(&t, "SELECT * FROM threads WHERE id = $1", id); err != nil {
		return store.Thread{}, fmt.Errorf("error getting thread: %w", err)
	}
	return t, nil
}

func (s *ThreadStore) CreateThread(t *store.Thread) error {
	if err := s.Get(t, "INSERT INTO threads VALUES ($1, $2, $3) RETURNING *",
		t.ID,
		t.Title,
		t.Description); err != nil {
		return fmt.Errorf("error creating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) UpdateThread(t *store.Thread) error {
	if err := s.Get(t, "UPDATE threads SET title = $1, description = $2 WHERE id = $3 RETURNING *",
		t.Title,
		t.Description,
		t.ID); err != nil {
		return fmt.Errorf("error updating thread: %w", err)
	}
	return nil
}

func (s *ThreadStore) DeleteThread(id uuid.UUID) error {
	if _, err := s.Exec("DELETE FROM threads WHERE id = $1", id); err != nil {
		return fmt.Errorf("error deleting thread: %w", err)
	}
	return nil
}
