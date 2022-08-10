package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tasks/models"
)

type TaskController struct{}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func (ctrl TaskController) AddTask(c *gin.Context) {
	var task models.Task
	c.BindJSON(&task)
	result, err := models.CreateTask(task)
	if strings.Contains(fmt.Sprint(err), "UNIQUE constraint") {
		arrayErr := strings.Split(fmt.Sprint(err), ":")
		taskErr := strings.Replace(arrayErr[1], ".", " ", -1)
		c.JSON(http.StatusConflict, gin.H{"message": taskErr + " already exists"})
		return
	} else if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create task"})
		return
	} else if result > 0 {
		c.JSON(http.StatusCreated, gin.H{
			"status":  "success",
			"message": "task created successfully",
		})
	}
}

func (ctrl TaskController) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	checkError(err)
	task, err := models.GetTaskById(i)
	checkError(err)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "task found",
			"data":    task,
		})
	}
}

func (ctrl TaskController) GetTasks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	pageInt, err := strconv.Atoi(page)
	checkError(err)
	limitInt, err := strconv.Atoi(limit)
	checkError(err)
	task, err := models.GetTasks(pageInt, limitInt)
	checkError(err)

	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No task found",
		})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Get List Task",
			"data":    task,
		})
	}
}

func (ctrl TaskController) UpdateTask(c *gin.Context) {
	var task models.Task
	err := c.BindJSON(&task)
	checkError(err)
	result, err := models.UpdateTask(task)
	checkError(err)

	if result > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "task updated successfully",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "task not found",
		})
	}
}

func (ctrl TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	checkError(err)
	result, err := models.DeleteTask(i)
	checkError(err)

	if result > 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "task deleted successfully",
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "task not found",
		})
	}
}

func (ctrl TaskController) MarkDone(c *gin.Context) {
	id := c.Param("id")
	i, err := strconv.Atoi(id)
	checkError(err)
	task, err := models.MarkDoneTask(i)
	checkError(err)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "task not found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "task updated successfully",
			"task":    task,
		})
	}
}
