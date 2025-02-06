package api

import (
	"Todo-Application/util"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

type TodoItemId struct {
	ID int `json:"id"`
}

type TodoRequest struct {
	Action string   `json:"action"`
	Item   TodoItem `json:"item"`
	ID     int      `json:"id"`
	Resp   chan TodoItem
}

var (
	todoItems = []TodoItem{}
	todoChan  = make(chan TodoRequest)
)

func Main() {
	loadTodosFromFile()
	go todoManager()

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
	var item TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON payload in create", err)
		return
	}

	todoChan <- TodoRequest{Action: "create", Item: item}
	writer.WriteHeader(http.StatusCreated)
}

func updateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in update", err)
		return
	}

	todoChan <- TodoRequest{Action: "update", Item: item}
	writer.WriteHeader(http.StatusOK)
}

func deleteTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item TodoItemId
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload for delete", err)
		return
	}

	todoChan <- TodoRequest{Action: "delete", ID: item.ID}
	writer.WriteHeader(http.StatusOK)
}

func getTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var item TodoItemId
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in get", err)
		return
	}

	respChan := make(chan TodoItem)
	todoChan <- TodoRequest{Action: "get", ID: item.ID, Resp: respChan}
	foundItem := <-respChan
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(foundItem)
}

func todoManager() {
	ctx := context.WithValue(context.Background(), "traceID", uuid.New())

	for req := range todoChan {
		switch req.Action {
		case "create":
			req.Item.ID = len(todoItems) + 1
			todoItems = append(todoItems, req.Item)
			saveTodosToFile(ctx)
		case "update":
			for i, item := range todoItems {
				if item.ID == req.Item.ID {
					todoItems[i].Desc = req.Item.Desc
					todoItems[i].Status = req.Item.Status
					saveTodosToFile(ctx)
					break
				}
			}
		case "delete":
			var newItems []TodoItem
			for _, item := range todoItems {
				if item.ID != req.ID {
					newItems = append(newItems, item)
				}
			}
			todoItems = newItems
			saveTodosToFile(ctx)
		case "get":
			var foundItem TodoItem
			for _, item := range todoItems {
				if item.ID == req.ID {
					foundItem = item
					break
				}
			}
			if req.Resp != nil {
				req.Resp <- foundItem
			}
		}
	}
}

func saveTodosToFile(ctx context.Context) {
	file, err := os.Create("list.json")
	if err != nil {
		util.LogError(ctx, "Error saving todos:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(todoItems)
	if (err != nil) {
		util.LogError(ctx, "Error encoding todos to file:", err)
	}
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