package main

import (
	"github.com/StevenSermeus/go-learning/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main () {
	godotenv.Load()

	server := gin.Default()	

	server.POST("/applications", controllers.CreateApplication)
	server.POST("/users", controllers.CreateUser)
	server.POST("/login", controllers.LoginHandler)
	server.Run()
}