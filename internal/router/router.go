package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Start() error {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowFiles = true
	config.AllowHeaders = []string{"*"}
	config.AllowMethods = []string{"*"}
	router.Use(cors.New(config))

	router.POST("/slice", SliceFile)

	return router.Run(":8080")
}
