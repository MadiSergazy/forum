package handler

import (
	"errors"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"html/template"
	"net/http"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/signin.html")
		if err != nil {
			h.l.Error.Printf("Parse file error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, temp); err != nil {
			h.l.Error.Printf("Execute tmpl error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		email := r.FormValue("email")
		password := r.FormValue("psw")
		user, err := h.service.GenerateToken(email, password)
		if err != nil {
			h.l.Error.Printf("Generate toke error: %s", err.Error())
			if errors.Is(err, service.InvalidData) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			http.SetCookie(w, &http.Cookie{
				Name:  "session_token",
				Value: user.Token,
				Path:  "/",
			})

			h.l.Info.Printf("the user is logged in6 email:%s", user.Login)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	default:
		h.l.Error.Printf("Method not allowed error:")
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		temp, err := template.ParseFiles("ui/signup.html")
		if err != nil {
			h.l.Error.Printf("Parse file error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err = temp.Execute(w, temp); err != nil {
			h.l.Error.Printf("Execute tmpl error:")
			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	case http.MethodPost:
		var signUp model.User
		email := r.FormValue("email")
		password := r.FormValue("psw")
		username := r.FormValue("username")
		repeatpass := r.FormValue("repeatspw")

		if repeatpass != password {
			h.l.Error.Println("Password not equal error:", password, repeatpass)
			errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		} else {
			signUp = model.User{
				Username: username,
				Login:    email,
				Password: password,
			}

		}

		if err := h.service.CreateUser(signUp); err != nil {
			h.l.Error.Printf("Create user error: %s", err.Error())
			if errors.Is(err, service.InvalidData) {
				errorHandler(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			errorHandler(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		h.l.Info.Printf("User created: email-%s", email)
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)

	default:
		h.l.Error.Println("Method not allowed error:")
		errorHandler(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:  "session_token",
		Value: "",
		Path:  "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
