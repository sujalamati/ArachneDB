package main

import (
	"fmt"
	"github.com/sujalamati/ArachneDB/api"
)

func main() {

	r := api.SetupRouter()
	// Run the server
	if err := r.Run(":8080"); err != nil {
		fmt.Println("Failed to start server:", err)
	}
	
}