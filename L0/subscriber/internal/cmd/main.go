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
			log.Println("Oh oh oh, ERROR:", e)
		}
	}()

	log.Println("Sub is running...(logging)")

	application, err := app.NewApp(".env")

	if err != nil {
		log.Fatalf("while creating app %v", err)
	}

	err = application.Run()
	if err != nil {
		log.Fatalf("while running app %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	sig := <-quit
	log.Printf("%v was called. Shutdown the app", sig)
}
