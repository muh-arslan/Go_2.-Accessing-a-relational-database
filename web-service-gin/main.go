package main

import (
	dataaccess "example/data-access"

	"github.com/gin-gonic/gin"
)

func main() {
	dataaccess.DBConnection()

	router := gin.Default()

	router.Run("localhost:8080")
}
