package todo

import (
	"Todo-Application/util"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"github.com/google/uuid"
)

type TodoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

var TodoItems = []TodoItem{}
var fileName = "list.json"

func Main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())
	
	// - When the application starts, load all to-do items from disk before adding new item
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		// create file if doesn't exist
		util.LogInfo(ctx, "Creating a new file!")
		file, err := os.Create(fileName)
		if err != nil {
			util.LogError(ctx, "Error creating file", err)
			return
		}
		defer file.Close()
	}

	// load file contents
	content, err := os.ReadFile(fileName)
	if err != nil {
		util.LogError(ctx, "Error reading file", err)
		return
	}

	if content != nil {
		json.Unmarshal(content, &TodoItems)
	}
	
	// - Create a command line application that uses flags to accept a to-do item adds it to an empty list of to-do items and prints the list to console
	var action string
	var id int
	var desc string
	var status string
	flag.StringVar(&action, "action", "add", "This is used to define what action should be performed")
	flag.IntVar(&id, "id", 0, "This is used to define what task should be completed")
	flag.StringVar(&desc, "desc", "empty will not be saved", "This is used to define what task should be completed")
	flag.StringVar(&status, "status", "notStarted", "This is the current status for the task")
	flag.Parse()

	var item = TodoItem{id, desc, status}
	fmt.Println(item)
	
	//Operations
	switch {
	case action == "add":
		addItem(item)
	case action == "remove":
		removeItem(TodoItems, item)
	case action == "update":
		updateItem(TodoItems, item)
	default:
		return
	}

	// - After printing the list of to-do items, save them to a file on disk
	jsonBytes, err := json.Marshal(TodoItems)
	if err != nil {
		util.LogError(ctx, "Error marshalling json", err)
		return
	}
	os.WriteFile(fileName, jsonBytes, 0644)
	util.LogInfo(ctx, "File saved successfully")
}

func addItem(item TodoItem) {
	TodoItems = append(TodoItems, item)
}

func removeItem(slice []TodoItem, item TodoItem) {
	result := []TodoItem{}
	for _, v := range slice {
		if v.ID != item.ID {
			result = append(result, v)
		}
	}
	TodoItems = result
}

func updateItem(result []TodoItem, item TodoItem) {
	for i, v := range result {
		if v.ID == item.ID {
			result[i].Desc = item.Desc
			result[i].Status = item.Status
		}
	}
	TodoItems = result;
}
