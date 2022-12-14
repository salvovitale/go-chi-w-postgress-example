package store

import (
	"github.com/google/uuid"
)

type Thread struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
}

type Post struct {
	ID            uuid.UUID `db:"id"`
	ThreadID      uuid.UUID `db:"thread_id"`
	Title         string    `db:"title"`
	Content       string    `db:"content"`
	Votes         int       `db:"votes"`
	CommentsCount int       `db:"comments_count"`
	ThreadTitle   string    `db:"thread_title"`
}

type Comment struct {
	ID      uuid.UUID `db:"id"`
	PostID  uuid.UUID `db:"post_id"`
	Content string    `db:"content"`
	Votes   int       `db:"votes"`
}

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`
	Password string    `db:"password"`
}

// Lets define what sort of storing and retrieving operations our database should be able to do on our entities

type ThreadStore interface {
	Threads() ([]Thread, error)
	Thread(id uuid.UUID) (Thread, error)
	CreateThread(t *Thread) error
	UpdateThread(t *Thread) error
	DeleteThread(id uuid.UUID) error
}

type PostStore interface {
	PostsByThread(threadID uuid.UUID) ([]Post, error)
	Posts() ([]Post, error)
	Post(id uuid.UUID) (Post, error)
	CreatePost(t *Post) error
	UpdatePost(t *Post) error
	DeletePost(id uuid.UUID) error
}

type CommentStore interface {
	CommentsByPost(postID uuid.UUID) ([]Comment, error)
	Comment(id uuid.UUID) (Comment, error)
	CreateComment(t *Comment) error
	UpdateComment(t *Comment) error
	DeleteComment(id uuid.UUID) error
}

type UserStore interface {
	User(id uuid.UUID) (User, error)
	UserByUsername(username string) (User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(id uuid.UUID) error
}

type Store interface {
	ThreadStore
	PostStore
	CommentStore
	UserStore
}
