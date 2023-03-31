// User REST API to perform CRUD operations in MongoDB
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

//	@title			User API
//	@version		1.0
//	@description	This is a User REST API to perform CRUD operations in MongoDB.

//	@contact.name	Anderson

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:3000
//	@BasePath	/api/v1
//	@schemes	http https

// @securityDefinitions.basic	BasicAuth
func main() {
	c := config.NewConfig()
	s := server.NewServer(c)

	c.Validate()

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
