package main

import (
	"caniteySnippetBox/internal/models"
	"caniteySnippetBox/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)


type snippetCreateForm struct {
	Title string `form:"title"`
	Content string `form:"content"`
	Expires int `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name string `form:"name"`
	Email string `form:"email"`
	Password string `form:"password"`
	validator.Validator `form:"-"`
}

func (a *application)home(w http.ResponseWriter, r *http.Request) {
	snippets, err := a.snippets.Latest()
	if err != nil {
		a.serverError(w, err)
		return
	}

	data := a.newTemplateData(r)
	data.Snippets = snippets
	a.render(w, http.StatusOK, "home.tmpl", data)
}

func (a *application)snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		a.notFound(w)
		return
	}

	snippet, err := a.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			a.notFound(w)
		} else {
			a.serverError(w, err)
		}
		return
	}


	data := a.newTemplateData(r)
	data.Snippet = snippet
	a.render(w, http.StatusOK, "view.tmpl", data)
}

func (a *application) snippetCreateForm(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	a.render(w, http.StatusOK, "create.tmpl", data)
}

func (a *application)snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	var form snippetCreateForm
	err = a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "this field cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "this field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "content cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must equal 1, 7 or 365")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := a.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		a.serverError(w, err)
		return
	}

	a.sessionManager.Put(r.Context(), "flash", "snippet created successfully")

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)

}


func (a *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := a.newTemplateData(r)
	data.Form = userSignupForm{}
	a.render(w, http.StatusOK, "signup.tmpl", data)
}

func (a *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm
	err := a.decodePostForm(r, &form)
	if err != nil {
		a.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be empty")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be empty")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Not a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be empty")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "this field must atleast be 8 characters long")

	if !form.Valid() {
		data := a.newTemplateData(r)
		data.Form = form
		a.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		return
	}

	fmt.Fprint(w, "created new user")
}

func (a *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user...")
}
func (a *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (a *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "logout user")
}
