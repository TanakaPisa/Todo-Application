package main

import (
	"Todo-Application/api"
	"Todo-Application/util"
	//"Todo-Application/web"
	"context"
	"github.com/google/uuid"
	//"Todo-Application/todo"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx := context.WithValue(context.Background(), "traceID", uuid.New())
	util.LogInfo(ctx, "Application started")
	// Create a channel to receive OS signals
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Run todo operations
	go api.Main()
	// todo.Main()
	// web.StaticPage()
	// web.DynamicPage()

	// Wait for a signal
	done := make(chan bool, 1)
	go func() {
		sig := <-channel 
		fmt.Printf("\nReceived signal: %s\n", sig)
		done <- true 
	}()
	<-done
}
