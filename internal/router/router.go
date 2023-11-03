package router

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Start() error {
	router := gin.Default()

	router.POST("/slice", SliceFile)

	return router.Run(":8080")
}
