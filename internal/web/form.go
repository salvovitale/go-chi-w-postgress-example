package web

import "encoding/gob"

func init() {
	// encode the form and errors so they can be sent to the client via the session
	gob.Register(CreatePostForm{})
	gob.Register(CreateThreadForm{})
	gob.Register(CreateCommentForm{})
	gob.Register(FormErrors{})
}

type FormErrors map[string]string

type CreatePostForm struct {
	Title   string
	Content string
	Errors  FormErrors
}

func (f *CreatePostForm) Validate() bool {
	f.Errors = make(FormErrors)

	if f.Title == "" {
		f.Errors["Title"] = "Title is required"
	}
	if f.Content == "" {
		f.Errors["Content"] = "Content is required"
	}

	return len(f.Errors) == 0
}

type CreateThreadForm struct {
	Title       string
	Description string
	Errors      FormErrors
}

func (f *CreateThreadForm) Validate() bool {
	f.Errors = make(FormErrors)

	if f.Title == "" {
		f.Errors["Title"] = "Title is required"
	}
	if f.Description == "" {
		f.Errors["Description"] = "Description is required"
	}

	return len(f.Errors) == 0
}

type CreateCommentForm struct {
	Content string
	Errors  FormErrors
}

func (f *CreateCommentForm) Validate() bool {
	f.Errors = make(FormErrors)

	if f.Content == "" {
		f.Errors["Content"] = "Content is required"
	}

	return len(f.Errors) == 0
}
