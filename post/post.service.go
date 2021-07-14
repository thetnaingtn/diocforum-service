package post

import (
	"context"
	"diocforum/middleware"
	"diocforum/user"
	"diocforum/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const postPath = "post"

func handlePostCreate(w http.ResponseWriter, r *http.Request) {
	var post Post
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "POST":
		err := json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		result, err := createPost(ctx, post)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		post.Id = result.InsertedID.(primitive.ObjectID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(post)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleEditPost(w http.ResponseWriter, r *http.Request) {
	var post Post
	objectId, _ := util.GetObjectIdFromUrl(r.URL.Path, fmt.Sprintf("/%s/edit", postPath))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	switch r.Method {
	case "PUT":
		if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		p, err := getPost(ctx, objectId)
		if post.UserId != p.UserId {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"User can only edit his own post"}`))
			return
		}

		_, err = updatePost(ctx, post, objectId)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		post.Id = objectId

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(post)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func handleDeletePost(w http.ResponseWriter, r *http.Request) {
	objectId, _ := util.GetObjectIdFromUrl(r.URL.Path, fmt.Sprintf("/%s/delete", postPath))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var ID user.Identity

	err := json.NewDecoder(r.Body).Decode(&ID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	p, _ := getPost(ctx, objectId)
	userId, _ := primitive.ObjectIDFromHex(ID.UserId)
	if p.UserId != userId {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error":"User can only delete his own post"}`))
		return
	}
	result, err := deletePost(ctx, objectId)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"deleted":{"count":"%d"}}`, result.DeletedCount)))
}

func SetupRoute(base string) {
	postCreateHandler := http.HandlerFunc(handlePostCreate)
	postEditHandler := http.HandlerFunc(handleEditPost)
	postDeleteHandler := http.HandlerFunc(handleDeletePost)

	http.Handle(fmt.Sprintf("%s/%s/new", base, postPath), middleware.CORS(middleware.Auth(postCreateHandler)))
	http.Handle(fmt.Sprintf("%s/%s/edit/", base, postPath), middleware.CORS(middleware.Auth(postEditHandler)))
	http.Handle(fmt.Sprintf("%s/%s/delete/", base, postPath), middleware.CORS(middleware.Auth(postDeleteHandler)))
}
