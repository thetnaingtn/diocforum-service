package auth

import (
	"context"
	"diocforum/middleware"
	"diocforum/session"
	"diocforum/user"
	"diocforum/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func handlerSignup(w http.ResponseWriter, r *http.Request) {
	var u user.User
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch r.Method {
	case "POST":
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		userId, err := user.CreateUser(ctx, u)

		if err != nil {
			if err == user.ErrUserExist {
				w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, user.ErrUserExist.Error())))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		u.Id = userId
		json.NewEncoder(w).Encode(u)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	var u user.User
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "POST":
		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user, err := user.GetUserByEmail(ctx, u.Email)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if util.Encrypt(u.Password) != user.Password {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"Password missmatch"}`))
			return
		}

		_, err = session.CreateSession(ctx, session.Session{UserId: user.Id, Email: user.Email})
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(user)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "POST":
		var ID user.Identity
		err := json.NewDecoder(r.Body).Decode(&ID)
		sess, err := session.GetSessionByUserId(ctx, ID.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = session.DeleteUserSession(ctx, sess.UserId)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"info":"user signout"}`))
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func SetUpRoute(apiBase string) {
	signupHandler := http.HandlerFunc(handlerSignup)
	loginHandler := http.HandlerFunc(handlerLogin)
	logoutHandler := http.HandlerFunc(handlerLogout)

	http.Handle(fmt.Sprintf("%s/signup", apiBase), middleware.CORS(signupHandler))
	http.Handle(fmt.Sprintf("%s/login", apiBase), middleware.CORS(loginHandler))
	http.Handle(fmt.Sprintf("%s/logout", apiBase), middleware.CORS(middleware.Auth(logoutHandler)))
}
