package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"userapi/config"
	"userapi/server"
)

func main() {
	c := config.NewConfig()
	s := server.NewServer(c)

	go func() {
		if err := s.Run(); err != nil {
			fmt.Println(fmt.Errorf("%v", err))
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("FATAL SERVER SHUTDOWN:", err)
	}

	fmt.Println("Server exiting")
}
