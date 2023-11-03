package main

import (
	"fmt"
	"github.com/AwesomeBFM/CloudSlicer/internal/router"
)

func main() {
	/*cmd := exec.Command("prusa-slicer", "--export-gcode", "--load", "./presets/pla.ini", "example.stl", "--output", "benchy.gcode", "--info")

	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Prusa Slicer: %v\n", err)
		return
	}

	outputStr := string(output)

	fmt.Printf("Output: %s\n", outputStr)*/

	apiRouter := router.NewRouter()

	err := apiRouter.Start()
	if err != nil {
		fmt.Printf("[ROUTER] [ERROR] Failed to start router! Error: %v\n", err)
		return
	}
}
