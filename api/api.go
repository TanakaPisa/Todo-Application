package api

import (
	"Todo-Application/concurrency"
	"Todo-Application/util"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func Main() {
	loadTodosFromFile()
	go concurrency.Main()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", createTodoHandler)
	mux.HandleFunc("POST /update", updateTodoHandler)
	mux.HandleFunc("POST /delete", deleteTodoHandler)
	mux.HandleFunc("POST /get", getTodoHandler)

	handler := util.CreateMiddleware(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	server.ListenAndServe()
}

func createTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item util.TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON payload in create", err)
		return
	}

	concurrency.TodoChan <- util.TodoRequest{Action: "create", Item: item}
	writer.WriteHeader(http.StatusCreated)
}

func updateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item util.TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in update", err)
		return
	}

	concurrency.TodoChan <- util.TodoRequest{Action: "update", Item: item}
	writer.WriteHeader(http.StatusOK)
}

func deleteTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item util.TodoItemId
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload for delete", err)
		return
	}

	concurrency.TodoChan <- util.TodoRequest{Action: "delete", ID: item.ID}
	writer.WriteHeader(http.StatusOK)
}

func getTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item util.TodoItemId
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in get", err)
		return
	}

	respChan := make(chan util.TodoItem)
	concurrency.TodoChan <- util.TodoRequest{Action: "get", ID: item.ID, Resp: respChan}
	foundItem := <-respChan
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(foundItem)
}

func loadTodosFromFile() {
	ctx := context.WithValue(context.Background(), "traceID", uuid.New())

	file, err := os.Open("list.json")
	if err != nil {
		if !os.IsNotExist(err) {
			util.LogError(ctx, "Error loading todos:", err)
		}
		return
	}
	defer file.Close()
}