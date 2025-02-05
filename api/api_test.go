package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func makeRequest(item TodoItem) *httptest.ResponseRecorder {
	// create items for test
	reqBody, _ := json.Marshal(item)
	req := httptest.NewRequest("POST", "/create", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Post
	rec := httptest.NewRecorder()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create", createTodoHandler)
	mux.ServeHTTP(rec, req)
	return rec
}


// Test two concurrent Create requests
func TestCreate(t *testing.T) {
	go todoManager()
	
	var wg sync.WaitGroup
	wg.Add(2)

	// First todo item
	go func() {
		defer wg.Done()
		item := TodoItem{ID: 1, Desc: "task1",Status: "pending"}
		rec := makeRequest(item)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected Status 201, got %d", rec.Code)
		}
	}()

	// Second todo item
	go func() {
		defer wg.Done()
		item := TodoItem{ID: 2, Desc: "task2",Status: "pending"}
		rec := makeRequest(item)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected Status 201, got %d", rec.Code)
		}
	}()

	wg.Wait() // Wait for both Goroutines to finish
}