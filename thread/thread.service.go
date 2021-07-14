package thread

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

const threadPath = "thread"

func handleThreadCreate(w http.ResponseWriter, r *http.Request) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "POST":
		var thread Thread
		err := json.NewDecoder(r.Body).Decode(&thread)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := createThread(ctx, thread)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		insertedThread, err := getThread(ctx, result.InsertedID.(primitive.ObjectID))
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(insertedThread)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// /thread/read/threadId
func handleGetThread(w http.ResponseWriter, r *http.Request) {
	var ctx, cancel = context.WithCancel(context.Background())
	objectId, err := util.GetObjectIdFromUrl(r.URL.Path, fmt.Sprintf("/%s/read", threadPath))
	defer cancel()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		thread, err := getThread(ctx, objectId)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(thread)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// /thread
func handleGetAllThreads(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "GET":
		threads, _ := getAllThreads(ctx)
		json.NewEncoder(w).Encode(threads)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// /thread/edit/threadId
func handleEditThread(w http.ResponseWriter, r *http.Request) {
	objectId, _ := util.GetObjectIdFromUrl(r.URL.Path, fmt.Sprintf("/%s/edit", threadPath))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	switch r.Method {
	case "PUT":
		var thread Thread
		if err := json.NewDecoder(r.Body).Decode(&thread); err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		t, err := getThread(ctx, objectId)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if t.UserId != thread.UserId {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"User can only edit his own thread"}`))
			return
		}

		_, err = updateThread(ctx, objectId, thread)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		thread.Id = objectId
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(thread)
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

// /thread/delete/threadId
func handleDeleteThread(w http.ResponseWriter, r *http.Request) {
	objectId, err := util.GetObjectIdFromUrl(r.URL.Path, fmt.Sprintf("/%s/delete", threadPath))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "DELETE":
		var Id user.Identity

		err := json.NewDecoder(r.Body).Decode(&Id)
		thread, err := getThread(ctx, objectId)

		if thread.UserId.Hex() != Id.UserId {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error":"User can only delete his own thread"}`))
			return
		}

		result, err := deleteThread(ctx, objectId)
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf(`{"deleted":{"count":"%d"}}`, result.DeletedCount)))
	case "OPTIONS":
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func SetupRoute(apiBase string) {
	allThreadsHandler := http.HandlerFunc(handleGetAllThreads)
	threadReadHandler := http.HandlerFunc(handleGetThread)
	threadEditHandler := http.HandlerFunc(handleEditThread)
	threadDeleteHandler := http.HandlerFunc(handleDeleteThread)
	threadCreateHandler := http.HandlerFunc(handleThreadCreate)

	http.Handle(fmt.Sprintf("%s/%s", apiBase, threadPath), middleware.CORS(allThreadsHandler))
	http.Handle(fmt.Sprintf("%s/%s/read/", apiBase, threadPath), middleware.CORS(threadReadHandler))
	http.Handle(fmt.Sprintf("%s/%s/new", apiBase, threadPath), middleware.CORS(middleware.Auth(threadCreateHandler)))
	http.Handle(fmt.Sprintf("%s/%s/edit/", apiBase, threadPath), middleware.CORS(middleware.Auth(threadEditHandler)))
	http.Handle(fmt.Sprintf("%s/%s/delete/", apiBase, threadPath), middleware.CORS(middleware.Auth(threadDeleteHandler)))
}
