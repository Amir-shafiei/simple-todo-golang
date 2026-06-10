package main

import (
	"awesomeProject19/db"
	"awesomeProject19/handlers"
	"awesomeProject19/middleware"
	"awesomeProject19/models"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var UserHandler handlers.UserHandler
	var TaskHandler handlers.TaskHandler
	database := db.Connect()
	err := database.AutoMigrate(&models.User{}, &models.Task{})
	if err != nil {
		return
	}
	UserHandler.DB = database
	TaskHandler.DB = database
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", nil)
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(200, "register.html", nil)
	})
	r.POST("/register", UserHandler.Register)
	r.POST("/login", UserHandler.Login)
	r.GET("/logout", UserHandler.Logout)
	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(200, "dashboard.html", nil)

	})
	api := r.Group("/api")
	api.POST("/tasks", middleware.Authenticate(), TaskHandler.AddTask)
	api.GET("/tasks", middleware.Authenticate(), TaskHandler.ListTasks)
	api.PUT("/tasks/", middleware.Authenticate(), TaskHandler.UpdateTask)
	api.DELETE("/tasks/:id", middleware.Authenticate(), TaskHandler.DeleteTask)
	api.PATCH("/tasks/:id/done", middleware.Authenticate(), TaskHandler.MarkAsDone)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run("0.0.0.0:" + port)
}
