package web

import "encoding/gob"

func init() {
	// encode the form and errors so they can be sent to the client via the session
	gob.Register(CreatePostForm{})
	gob.Register(CreateThreadForm{})
	gob.Register(CreateCommentForm{})
	gob.Register(RegisterForm{})
	gob.Register(LoginForm{})
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

type RegisterForm struct {
	Username      string
	Password      string
	UsernameTaken bool

	Errors FormErrors
}

func (f *RegisterForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Username == "" {
		f.Errors["Username"] = "Please enter a username."
	} else if f.UsernameTaken {
		f.Errors["Username"] = "This username is already taken."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password."
	} else if len(f.Password) < 8 {
		f.Errors["Password"] = "Your password must be at least 8 characters long."
	}

	return len(f.Errors) == 0
}

type LoginForm struct {
	Username             string
	Password             string
	IncorrectCredentials bool

	Errors FormErrors
}

func (f *LoginForm) Validate() bool {
	f.Errors = FormErrors{}

	if f.Username == "" {
		f.Errors["Username"] = "Please enter a username."
	} else if f.IncorrectCredentials {
		f.Errors["Username"] = "Username or password is incorrect."
	}

	if f.Password == "" {
		f.Errors["Password"] = "Please enter a password."
	}

	return len(f.Errors) == 0
}
