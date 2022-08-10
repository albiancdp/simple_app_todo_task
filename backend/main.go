package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tasks/controllers"
	"github.com/tasks/models"
)

var TaskController = new(controllers.TaskController)

func main() {
	r := gin.Default()
	r.Use(CORSMiddleware())
	err := models.ConnectDatabase()
	checkError(err)
	errMigrate := models.Migration()
	checkError(errMigrate)

	v1 := r.Group("/api/v1")
	{
		v1.POST("task", TaskController.AddTask)
		v1.GET("task-list", TaskController.GetTasks)
		v1.GET("task/:id", TaskController.GetTaskById)
		v1.POST("task/update", TaskController.UpdateTask)
		v1.GET("task/delete/:id", TaskController.DeleteTask)
		v1.GET("task/mark-done/:id", TaskController.MarkDone)
	}
	r.Run()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, HEAD, PATCH, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
