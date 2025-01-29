package api

import (
	"Todo-Application/util"
	"context"
	"encoding/json"
	"fmt"
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

var todoItems = []TodoItem{}

var ctx context.Context

func Main() {
	loadTodosFromFile()
	mux := http.NewServeMux()
	mux.HandleFunc("/create", createTodoHandler)
	mux.HandleFunc("/update", updateTodoHandler)
	mux.HandleFunc("/delete", deleteTodoHandler)
	mux.HandleFunc("/get", getTodoHandler)
	mux.HandleFunc("/list", ListTodoHandler)

	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())
	handler := util.CreateMiddleware(ctx, mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	server.ListenAndServe()
}

func createTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var item TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, "Invalid JSON payload in create", http.StatusBadRequest)
		return
	}

	item.ID = len(todoItems) + 1
	todoItems = append(todoItems, item)
	saveTodosToFile()

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(item)
}

func updateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var updateItem TodoItem
	err := json.NewDecoder(request.Body).Decode(&updateItem)
	if err != nil {
		http.Error(writer, "Invalid JSON Payload in update", http.StatusBadRequest)
		return
	}

	for i, item := range todoItems {
		if item.ID == updateItem.ID {
			todoItems[i].Desc = updateItem.Desc
			todoItems[i].Status = updateItem.Status
			saveTodosToFile()

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(updateItem)
			return
		}
	}

	http.Error(writer, "Todo Item not found", http.StatusNotFound)
}

func deleteTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method in delete", http.StatusMethodNotAllowed)
		return
	}

	var removeItem TodoItemId
	result := []TodoItem{}
	err := json.NewDecoder(request.Body).Decode(&removeItem)
	if err != nil {
		http.Error(writer, "Invalid JSON Payload", http.StatusBadRequest)
		return
	}

	for _, item := range todoItems {
		fmt.Println("item:", item.ID)
		fmt.Println("item:", removeItem.ID)
		if item.ID != removeItem.ID {
			result = append(result, item)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(item)
		}
	}

	todoItems = result
	saveTodosToFile()
}

func getTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var getItem TodoItemId
	err := json.NewDecoder(request.Body).Decode(&getItem)
	if err != nil {
		http.Error(writer, "Invalid JSON Payload in get", http.StatusBadRequest)
		return
	}

	for _, item := range todoItems {
		if item.ID == getItem.ID {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(item)
			return
		}
	}

	http.Error(writer, "Todo Item not found", http.StatusNotFound)
}

func ListTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	for _, item := range todoItems {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(item)
	}
}

func saveTodosToFile() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())
	
	file, err := os.Create("list.json")
	if err != nil {
		util.LogError(ctx, "Error saving todos:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(todoItems); err != nil {
		util.LogError(ctx, "Error encoding todos to file:", err)
	}
}

func loadTodosFromFile() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())

	file, err := os.Open("list.json")
	if err != nil {
		if !os.IsNotExist(err) {
			util.LogError(ctx, "Error loading todos:", err)
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todoItems); err != nil {
		util.LogError(ctx,"Error decoding todos from file:", err)
	}
}