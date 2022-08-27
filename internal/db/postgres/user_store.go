package postgres

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/salvovitale/go-chi-w-postgress-example/internal/db/store"
)

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{DB: db}
}

type UserStore struct {
	*sqlx.DB
}

func (s *UserStore) User(id uuid.UUID) (store.User, error) {
	var u store.User
	if err := s.Get(&u, `SELECT * FROM users WHERE id = $1`, id); err != nil {
		return store.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) UserByUsername(username string) (store.User, error) {
	var u store.User
	if err := s.Get(&u, `SELECT * FROM users WHERE username = $1`, username); err != nil {
		return store.User{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *UserStore) Users() ([]store.User, error) {
	var uu []store.User
	if err := s.Select(&uu, `SELECT * FROM users`); err != nil {
		return []store.User{}, fmt.Errorf("error getting users: %w", err)
	}
	return uu, nil
}

func (s *UserStore) CreateUser(u *store.User) error {
	if err := s.Get(u, `INSERT INTO users VALUES ($1, $2, $3) RETURNING *`,
		u.ID,
		u.Username,
		u.Password); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	return nil
}

func (s *UserStore) UpdateUser(u *store.User) error {
	if err := s.Get(u, `UPDATE users SET username = $1, password = $2 WHERE id = $3 RETURNING *`,
		u.Username,
		u.Password,
		u.ID); err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	return nil
}

func (s *UserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec(`DELETE FROM users WHERE id = $1`, id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
