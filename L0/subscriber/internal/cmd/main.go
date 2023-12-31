package main

import (
	"log"
	"ordermngmt/internal/app"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer func() {
		e := recover()
		if e != nil {
			log.Println("error was recovered", e)
		}
	}()

	log.Println("Sub is running...")

	application, err := app.NewApp(".env")

	if err != nil {
		log.Fatalf("while creating app %v", err)
	}

	errChan := make(chan error, 1)
	application.Run(errChan)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-errChan:
		log.Fatalf("while running app %v", err)
	case sig := <-quit:
		log.Printf("%v was called. Shutdown the app", sig)
	}

}
