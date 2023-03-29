package server

import (
	"context"
	"fmt"
	"net/http"
	"userapi/config"
	"userapi/users"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Server interface {
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	config config.Config
	srv    *http.Server
}

func NewServer(config config.Config) Server {
	return &server{
		config: config,
	}
}

func (s *server) Run() error {
	router := gin.Default()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router,
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(s.config.DBURI))
	if err != nil {
		panic(err)
	}

	api := router.Group("/api")
	users.AddRoutes(api, s.config, client)

	fmt.Printf("Running server on port : %v \n", s.config.Port)
	return s.srv.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
