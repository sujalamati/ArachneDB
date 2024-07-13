package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the Gin router with all the routes
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/username", createUserHandler)

	r.GET("/user/get/collection/:username/:collection_name", getColletionHandler)
	r.POST("/user/create/collection", createCollectionHandler)
	r.POST("/user/delete/collection", deleteCollectionHandler)

	r.GET("/user/get/key/:username/:collection/:key", getKeyHandler)
	r.POST("/user/delete/key", deleteKeyHandler)
	r.POST("/user/create/key", createKeyHandler)

	return r
}