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
 
var (
	todoItems = []TodoItem{}
	createChan = make(chan TodoItem)
	updateChan = make(chan TodoItem)
	deleteChan = make(chan int)
	getChan = make(chan int)
	getResp = make(chan TodoItem)
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

	createChan <- item
	writer.WriteHeader(http.StatusCreated)
}

func updateTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var updateItem TodoItem
	err := json.NewDecoder(request.Body).Decode(&updateItem)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in update", err)
		return
	}

	updateChan <- updateItem
	writer.WriteHeader(http.StatusOK)
}

func deleteTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var removeItem TodoItemId
	err := json.NewDecoder(request.Body).Decode(&removeItem)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload for delete", err)
		return
	}

	deleteChan <- removeItem.ID
	writer.WriteHeader(http.StatusOK)
}

func getTodoHandler(writer http.ResponseWriter, request *http.Request) {
	var getItem TodoItemId
	err := json.NewDecoder(request.Body).Decode(&getItem)
	if err != nil {
		util.LogError(request.Context(), "Invalid JSON Payload in get", err)
		return
	}

	getChan <- getItem.ID
	item := <- getResp
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(item)

	//util.LogError(request.Context(), "Todo not found", err)
}

func todoManager() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())

	for {
		select {
		case item := <- createChan:
			item.ID = len(todoItems) + 1
			todoItems = append(todoItems, TodoItem{ID: item.ID, Desc: item.Desc, Status: item.Status})
			saveTodosToFile(ctx)
		case update := <-updateChan:
			for i, item := range todoItems {
				if item.ID == update.ID {
					todoItems[i].Desc = update.Desc
					todoItems[i].Status = update.Status
					saveTodosToFile(ctx)
					break
				}
			}
		case id := <-deleteChan:
			var newItems []TodoItem
			for _, item := range todoItems {
				if item.ID != id {
					newItems = append(newItems, item)
				}
			}
			todoItems = newItems
			saveTodosToFile(ctx)

		case id := <-getChan:
			var foundItem TodoItem
			for _, item := range todoItems {
				if item.ID == id {
					foundItem = item
					break
				}
			}
			getResp <- foundItem
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
}