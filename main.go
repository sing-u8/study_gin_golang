package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Recipe struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"publishedAt"`
}

func IndexHandler(c *gin.Context) {
	title := c.Params.ByName("title")
	c.JSON(200, gin.H{
		"message": "hell " + title,
	})
}

func main() {
	router := gin.Default()
	router.GET("/:title", IndexHandler)
	router.Run()
}
