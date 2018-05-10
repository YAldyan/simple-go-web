package controller

import (
	"example-go-web/util"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleUserNew(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Display Home Page
	RenderTemplate(w, r, "users/new", nil)
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user, err := NewUser(
		r.FormValue("username"),
		r.FormValue("email"),
		r.FormValue("password"),
	)

	if err != nil {
		if util.IsValidationError(err) {
			RenderTemplate(w, r, "users/new", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
		return
	}

	/*
		menyimpan user yg telah d create ke dalam
		user list di file users.json,
	*/
	err = GlobalUserStore.Save(user)
	if err != nil {
		panic(err)
		return
	}

	// Create a new session
	session := NewSession(w)
	session.UserID = user.ID

	/*
		menyimpan session yg telah d create ke dalam
		session list di file users.json,
	*/
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/?flash=User+created", http.StatusFound)
}

func HandleUserEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	user := RequestUser(r)
	RenderTemplate(w, r, "users/edit", map[string]interface{}{
		"User": user,
	})
}

func HandleUserUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(r)
	email := r.FormValue("email")
	currentPassword := r.FormValue("currentPassword")
	newPassword := r.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	if err != nil {
		if util.IsValidationError(err) {
			RenderTemplate(w, r, "users/edit", map[string]interface{}{
				"Error": err.Error(),
				"User":  user,
			})
			return
		}
		panic(err)
	}

	err = GlobalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/account?flash=User+updated", http.StatusFound)
}

func HandleUserShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := GlobalUserStore.Find(params.ByName("userID"))
	if err != nil {
		panic(err)
	}

	// 404
	if user == nil {
		http.NotFound(w, r)
		return
	}

	images, err := GlobalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	RenderTemplate(w, r, "users/show", map[string]interface{}{
		"Images": images,
		"User":   user,
	})
}
