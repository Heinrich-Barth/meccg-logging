package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JsonError(message string) map[string]string {
	return map[string]string{
		"message": message,
	}
}

func ErrorHandler(c *gin.Context) {

	// first, handle the actual route
	c.Next()

	// no errors, so all is fine
	if len(c.Errors) == 0 {
		return
	}

	// return error message
	var err = c.Errors.Last().Err
	var message = err.Error()

	c.JSON(http.StatusInternalServerError, JsonError(message))
	log.Println(message)
}

func initServer() {

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		ErrorHandler(c)
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusNoContent, nil)
	})

	router.GET("/games", func(c *gin.Context) {
		var list = listGameFiles()
		if list == nil {
			c.JSON(http.StatusOK, []string{})
		} else {
			c.JSON(http.StatusOK, list)
		}
	})

	router.Static("/games", "./games")

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.Run()
}
