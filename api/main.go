package main

import (
  "net/http"
  "github.com/gin-gonic/gin"
)

func main() {
  r := gin.Default()
  r.GET("/", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "hello": "world",
    })
  })
  r.NoRoute(func(c *gin.Context) {
    c.JSON(http.StatusNotFound, gin.H{
      "error": "not found",
    })
  })
  r.Run()
}
