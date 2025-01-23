package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
)

type todoItem struct{
	ID int `json:"id"`
	Desc string `json:"desc"`
	Status string `json:"status"`
}

var todoItems = []todoItem{}
var fileName = "list.json"

func main() {
	// - When the application starts, load all to-do items from disk before adding new item
	_, err := os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		// create file if doesn't exist	
		fmt.Println("Creating a new file!")
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} 
	
	// load file contents
	content, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	
	if content != nil {
		json.Unmarshal(content, &todoItems)
	}
	

	// - Allow the user to update the description of a to-do item or delete it
	// remove item
	var removeId int
	flag.IntVar(&removeId,"removeId", 0, "This is used to remove task from list by ID")

	//update item
	var updateId int
	var updateDesc string
	var updateStatus string
	flag.IntVar(&updateId,"updateId", 0, "This is the Id for todo item that requires update")
	flag.StringVar(&updateDesc,"updateDesc", "updated", "Update task description")
	flag.StringVar(&updateStatus,"updateStatus", "not started", "Update task for status")

	// - Create a command line application that uses flags to accept a to-do item adds it to an empty list of to-do items and prints the list to console
	// Add item 
	var id int
	var desc string
	var status string
	flag.IntVar(&id,"id", 0, "This is used to define what task should be completed")
	flag.StringVar(&desc,"desc", "", "This is used to define what task should be completed")
	flag.StringVar(&status,"status", "notStarted", "This is the current status for the task")
	flag.Parse()
	
	var item = todoItem{id, desc, status}
	fmt.Println(item)

	//Operations
	if item.Desc != ""  {
		todoItems = append(todoItems, item)
	}
	
	if removeId > 0 {
		todoItems = removeTodoItem(todoItems, removeId)
	}

	if updateId > 0 {
		updateTodoItem(todoItems, updateId, updateDesc, updateStatus)
	}

	// - After printing the list of to-do items, save them to a file on disk
	jsonBytes, err := json.Marshal(todoItems)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))
	os.WriteFile(fileName, jsonBytes, 0644)
}

func removeTodoItem(slice []todoItem, target int)  []todoItem {
	result := []todoItem{}
	for i, v := range slice {
		fmt.Println(i)
		if v.ID != target {
			result = append(result, v)
		}
	}
	return result
}

func updateTodoItem(slice []todoItem, target int,desc string, status string)  []todoItem {
	for i, v := range slice {
		if v.ID == target {
			slice[i].Desc = desc
			slice[i].Status = status
		}
	}
	return slice
}