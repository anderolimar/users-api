// This package contains the server to run User API
package server

import (
	"context"
	"fmt"
	"net/http"
	"userapi/config"
	"userapi/users"

	_ "userapi/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

// Returns a new instance of Server
func NewServer(config config.Config) Server {
	return &server{
		config: config,
	}
}

// Method to run the server
func (s *server) Run() error {
	router := gin.Default()

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.config.Port),
		Handler: router,
	}

	// Rate Limiter
	var limiter = NewIPRateLimiter(s.config.RateLimit, s.config.RateLimitTokens)

	// MongoDB client connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(s.config.DBURI))
	if err != nil {
		panic(err)
	}

	// CORS
	router.Use(Cors)

	apiV1 := router.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		s.config.ApiUser: s.config.ApiPass,
	}))

	apiV1.Use(limitMiddleware(limiter))
	users.AddRoutes(apiV1, s.config, client)

	// API Documentation with swagger
	router.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Printf("Running server on port : %v \n", s.config.Port)
	return s.srv.ListenAndServe()
}

// Method to shutdown the server
func (s *server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
