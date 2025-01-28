package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

var todoItems = []TodoItem{}

func Main() {
	loadTodosFromFile()
	mux := http.NewServeMux()
	mux.HandleFunc("/create", createTodoHandler)
	mux.HandleFunc("/update", updateTodoHandler)
	mux.HandleFunc("/delete", deleteTodoHandler)
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
		http.Error(writer, "Invalid JSON Payload", http.StatusBadRequest)
		return
	}

	for i, item := range todoItems {
		if item.ID == updateItem.ID {
			todoItems[i].Desc = updateItem.Desc
			todoItems[i].Status = updateItem.Status
			saveTodosToFile()

			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			json.NewEncoder(writer).Encode(item)
			return
		}
	}

	http.Error(writer, "Todo Item not found", http.StatusNotFound)
}

func deleteTodoHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(writer, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var removeItem TodoItem
	result := []TodoItem{}
	err := json.NewDecoder(request.Body).Decode(&removeItem)
	if err != nil {
		http.Error(writer, "Invalid JSON Payload", http.StatusBadRequest)
		return
	}

	for _, item := range todoItems {
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

func saveTodosToFile() {
	file, err := os.Create("list.json")
	if err != nil {
		fmt.Println("Error saving todos:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(todoItems); err != nil {
		fmt.Println("Error encoding todos to file:", err)
	}
}

func loadTodosFromFile() {
	file, err := os.Open("list.json")
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error loading todos:", err)
		}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&todoItems); err != nil {
		fmt.Println("Error decoding todos from file:", err)
	}
}
