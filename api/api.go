package api

import (
	"encoding/json"
	"net/http"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

var todoItems = []TodoItem{}

func Main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", createTodoHandler)
	http.ListenAndServe(":8080", mux)
}

func createTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	
	var item TodoItem
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	item.ID = len(todoItems) + 1
	todoItems = append(todoItems, item)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(item)
}