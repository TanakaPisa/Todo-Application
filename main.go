package main

import (
	"Todo-Application/util"
	"context"
	"github.com/google/uuid"
	"Todo-Application/api"
	//"Todo-Application/todo"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", uuid.New())
	util.LogInfo(ctx, "Application started")
	// Create a channel to receive OS signals
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)

	// Run todo operations
	//todo.Main()
	go api.Main()

	// Wait for a signal
	done := make(chan bool, 1)
	go func() {
		sig := <-channel 
		fmt.Printf("\nReceived signal: %s\n", sig)
		done <- true 
	}()
	<-done
}
