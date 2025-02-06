package concurrency

import (
	"Todo-Application/util"
	"context"
	"encoding/json"
	"os"
	"github.com/google/uuid"
)

var (
	todoItems = []util.TodoItem{}
	TodoChan  = make(chan util.TodoRequest)
)

func Main() {
	go todoManager()
}

func todoManager() {
	ctx := context.WithValue(context.Background(), "traceID", uuid.New())

	for req := range TodoChan {
		switch req.Action {
		case "create":
			createTodo(ctx, req.Item)
		case "update":
			updateTodo(ctx, req.Item)
		case "delete":
			deleteTodo(ctx, req.ID)
		case "get":
			getTodos(req)
		}
	}
}

func createTodo(ctx context.Context, item util.TodoItem) {
	item.ID = len(todoItems) + 1
	todoItems = append(todoItems, item)
	saveTodosToFile(ctx)
}

func updateTodo(ctx context.Context, item util.TodoItem) {
	for i, t := range todoItems {
		if t.ID == item.ID {
			todoItems[i].Desc = item.Desc
			todoItems[i].Status = item.Status
			saveTodosToFile(ctx)
			break
		}
	}
}

func deleteTodo(ctx context.Context, id int) {
	var newItems []util.TodoItem
	for _, item := range todoItems {
		if item.ID != id {
			newItems = append(newItems, item)
		}
	}
	todoItems = newItems
	saveTodosToFile(ctx)
}

func getTodos(req util.TodoRequest) {
	var foundItem util.TodoItem
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