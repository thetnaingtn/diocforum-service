package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
)

func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		switch r.Method {
		case "POST", "GET", "PUT", "DELETE":
			var ID user.Identity
			body, _ := ioutil.ReadAll(r.Body)

			if err := json.Unmarshal(body, &ID); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if ID.UserId == "" {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error":"User identity must be present"}`))
				return
			}
			_, err := session.GetSessionByUserId(ctx, ID.UserId)
			if err == mongo.ErrNoDocuments {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"error":"User doesn't sign in"}`))
				return
			}
			r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			handler.ServeHTTP(w, r)
		case "OPTIONS":
			return
		}
	})
}
