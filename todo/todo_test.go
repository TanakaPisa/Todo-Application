package todo

import (
	"os"
	"testing"
)

func TestAddTodo(t *testing.T) {
	item := TodoItem{ID: 1, Desc : "write unit tests", Status: "started"}
	TodoItems = append(TodoItems, item)

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "testfile_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	Main()

	// Read the file to validate the content
	data, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to read the file: %v", err)
	}

	// Verify the content
	content := ""
	if string(data) != content {
		t.Errorf("Expected content %q, but got %q", content, string(data))
	}
}