package todo

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"
)

type todoItem struct {
	ID     int    `json:"id"`
	Desc   string `json:"desc"`
	Status string `json:"status"`
}

var todoItems = []todoItem{}
var fileName = "list.json"

type contextKey string
const traceIdKey contextKey = "traceID"

func Main() {
	// create context
	ctx := context.Background()
	ctx = context.WithValue(ctx, traceIdKey, fmt.Sprintf("trace-%d", time.Now().UnixNano()))

	// Use the log package to log errors and information to the console
	traceIdVal, ok := ctx.Value(traceIdKey).(string)
	if !ok {
		traceIdVal = "unknown" 
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})).With("traceID", traceIdVal)
	
	// - When the application starts, load all to-do items from disk before adding new item
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		// create file if doesn't exist
		fmt.Println("Creating a new file!")
		file, err := os.Create(fileName)
		if err != nil {
			logger.Error("Error creating file", slog.String("file", fileName), slog.Any("error", err))
			return
		}
		defer file.Close()
	}

	// load file contents
	content, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error("Error reading file", slog.String("file", fileName), slog.Any("error", err))
		return
	}

	if content != nil {
		json.Unmarshal(content, &todoItems)
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

	var item = todoItem{id, desc, status}
	fmt.Println(item)

	//Operations
	switch {
	case action == "add":
		addItem(item)
	case action == "remove":
		removeItem(todoItems, item)
	case action == "update":
		updateItem(todoItems, item)
	default:
		return
	}

	// - After printing the list of to-do items, save them to a file on disk
	jsonBytes, err := json.Marshal(todoItems)
	if err != nil {
		logger.Error("Error marshalling json", slog.String("file", fileName), slog.Any("error", err))
		return
	}
	os.WriteFile(fileName, jsonBytes, 0644)
	logger.Info("File saved successfully", slog.String("file", fileName))
}

func addItem(item todoItem) {
	todoItems = append(todoItems, item)
}

func removeItem(slice []todoItem, item todoItem) {
	result := []todoItem{}
	for _, v := range slice {
		if v.ID != item.ID {
			result = append(result, v)
		}
	}
	todoItems = result
}

func updateItem(result []todoItem, item todoItem) {
	for i, v := range result {
		if v.ID == item.ID {
			result[i].Desc = item.Desc
			result[i].Status = item.Status
		}
	}
	todoItems = result;
}
