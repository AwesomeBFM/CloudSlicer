package main

import (
	"fmt"
	"github.com/AwesomeBFM/CloudSlicer/internal/router"
)

func main() {
	apiRouter := router.NewRouter()

	err := apiRouter.Start()
	if err != nil {
		fmt.Printf("[ROUTER] [ERROR] Failed to start router! Error: %v\n", err)
		return
	}
}
