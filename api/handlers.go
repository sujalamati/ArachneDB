package api

import(
	"github.com/gin-gonic/gin"
	"github.com/sujalamati/ArachneDB/pkg"
	"net/http"
	"fmt"
)

func createUserHandler(c *gin.Context) {
	// Extract username from the header
	var json struct {
        Username string `json:"username"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	username := json.Username
	
	// Initialize the database
	filename := username + ".adb"

	db, err := ArachneDB.Open(filename, ArachneDB.DefaultOptions)

	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success message
	c.JSON(http.StatusOK, fmt.Sprintf("Username %s processed successfully", username))
}

func getColletionHandler(c *gin.Context) {
	// Extract username and collection name from parameters
	username := c.Param("username")
	collectionName := c.Param("collection_name")

	// Initialize the database
	filename := username + ".adb"
	db, err := ArachneDB.Open(filename, ArachneDB.DefaultOptions)
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	

	tx := db.ReadTx()
	// Get the collection
	col, err := tx.GetCollection([]byte(collectionName))
	if err != nil{
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	if col ==nil{
		c.JSON(http.StatusInternalServerError, "Collection not found!")
		return
	}

	// Respond with collection data and success message
	c.JSON(http.StatusCreated, fmt.Sprintf("Entered %s successfully", collectionName))
}

func createCollectionHandler(c *gin.Context) {
	var json struct {
        Username string `json:"username"`
		Collection string `json:"collection_name"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	username := json.Username
	collectionName := json.Collection

	// Initialize the database
	filename := username + ".adb"
	db, err := ArachneDB.Open(filename, ArachneDB.DefaultOptions)

	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx := db.WriteTx()
	// Create the collection
	_, err = tx.CreateCollection([]byte(collectionName))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()

	// Respond with success message
	c.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Collection %s created successfully", collectionName)})
}

func deleteCollectionHandler(c *gin.Context) {

	var json struct {
        Username string `json:"username"`
		Collection string `json:"collection_name"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	username := json.Username
	collectionName := json.Collection
	

	// Initialize the database
	filename := username + ".adb"
	db, err := ArachneDB.Open(filename, ArachneDB.DefaultOptions)

	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := db.WriteTx()
	// Delete the collection
	err = tx.DeleteCollection([]byte(collectionName))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	
	// Respond with success message
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Collection %s deleted successfully", collectionName)})
}

func getKeyHandler(c *gin.Context) {
	username := c.Param("username")
	collection := c.Param("collection")
	key := c.Param("key")

	db, err := ArachneDB.Open(username+".adb", ArachneDB.DefaultOptions)
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := db.ReadTx()

	// Perform operations with the database based on username, collection, and key
	// For example:
	c1,err:=tx.GetCollection([]byte(collection))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	i, err := c1.Find([]byte(key))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	if i == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	// Respond with the key and value
	c.JSON(http.StatusOK, fmt.Sprintf("Key is :%s, Value is :%s",i.Key(),i.Value()))

}

func createKeyHandler(c *gin.Context) {
	var json struct {
        Username string `json:"username"`
		Collection string `json:"collection"`
		Key string `json:"key"`
		Value string `json:"value"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	
	username := json.Username
	collection := json.Collection
	key := json.Key
	value:=json.Value

	db, err := ArachneDB.Open(username+".adb", ArachneDB.DefaultOptions)
	defer func() {
		_ = db.Close()
	}()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := db.WriteTx()

	c1,err:=tx.GetCollection([]byte(collection))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = c1.Put([]byte(key),[]byte(value))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Key %s created successfully", key)})
}

func deleteKeyHandler(c *gin.Context) {
	var json struct {
        Username string `json:"username"`
		Collection string `json:"collection"`
		Key string `json:"key"`
    }

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	
	username := json.Username
	collection := json.Collection
	key := json.Key

	


	db, err := ArachneDB.Open(username+".adb", ArachneDB.DefaultOptions)
	defer func() {
		_ = db.Close()
	}()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tx := db.WriteTx()

	// Perform operations with the database based on username, collection, and key
	// For example:
	c1,err:=tx.GetCollection([]byte(collection))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = c1.Remove([]byte(key))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Key %s deleted successfully", key)})
}