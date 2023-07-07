package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhayt/clean-arch-tmp-forum/internal/model"
	"github.com/zhayt/clean-arch-tmp-forum/internal/service"
	"github.com/zhayt/clean-arch-tmp-forum/logger"
	"net/http"
	"time"
)

type Middleware struct {
	user service.Authorization
	l    *logger.Logger
}

func (m *Middleware) LogRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.l.Info.Println(fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.String()))

		next.ServeHTTP(w, r)
	}
}

func NewMiddleware(user service.Authorization, l *logger.Logger) *Middleware {
	return &Middleware{user: user, l: l}
}

func (m *Middleware) WithSessionBlocked(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				next.ServeHTTP(w, r)
				return
			}

			m.l.Error.Printf("failed to get cookie: %s", err.Error())
			return
		}

		user, err := m.user.GetUserByToken(token.Value)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if user.TokenDuration.Before(time.Now()) {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (m *Middleware) AuthorizationRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user model.User
		token, err := r.Cookie("session_token")
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, model.User{})))
			case errors.Is(err, token.Valid()):
				m.l.Error.Println("invalid cookie value")
			}
			m.l.Error.Println("failed to get cookie")
			return
		}

		user, err = m.user.GetUserByToken(token.Value)
		m.l.Info.Println("token:", token.Value)
		if err != nil {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, model.User{})))
			return
		}

		m.l.Info.Println("User with token", user.Login, token.Value)
		if user.TokenDuration.Before(time.Now()) {
			if err := m.user.DeleteToken(user.Token); err != nil {
				fmt.Println(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), model.CtxUserKey, user)))
	}
}
