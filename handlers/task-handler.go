package handlers

import (
	"awesomeProject19/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	DB *gorm.DB
}

func (t *TaskHandler) AddTask(c *gin.Context) {
	var task models.Task
	var existingTask models.Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	task.UserID = userID
	result := t.DB.Where("title = ? AND user_id = ?",
		task.Title,
		userID).First(&existingTask)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task already exists"})
		return
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	t.DB.Create(&task)
	c.JSON(http.StatusOK, gin.H{"data": task})

}
func (t *TaskHandler) UpdateTask(c *gin.Context) {
	var task models.Task
	var existingTask models.Task
	err := c.ShouldBindJSON(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))
	result := t.DB.Where("id = ? AND user_id = ?", task.ID, userID).First(&existingTask)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}
	existingTask.Title = task.Title
	t.DB.Save(&existingTask)
	c.JSON(http.StatusOK, gin.H{"message": "task updated"})

}
func (t *TaskHandler) DeleteTask(c *gin.Context) {
	var task models.Task
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))
	result := t.DB.Where("id = ? AND user_id= ?", c.Param("id"), userID).First(&task)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	t.DB.Delete(&task)
	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})

}
func (t *TaskHandler) ListTasks(c *gin.Context) {
	var tasks []models.Task

	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	statusParam := c.Query("status")

	pageInt, err1 := strconv.Atoi(page)
	limitInt, err2 := strconv.Atoi(limit)

	if err1 != nil || err2 != nil || pageInt < 1 || limitInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination"})
		return
	}
	offset := (pageInt - 1) * limitInt
	query := t.DB.Where("user_id = ?", userID)

	if statusParam != "" {
		status := statusParam == "true"
		query = query.Where("status = ?", status)
	}

	result := query.
		Limit(limitInt).
		Offset(offset).
		Find(&tasks)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"page":  pageInt,
		"limit": limitInt,
		"data":  tasks,
	})
}
func (t *TaskHandler) MarkAsDone(c *gin.Context) {
	var task models.Task
	userIDVal, _ := c.Get("user_id")
	userID := uint(userIDVal.(float64))
	taskID := c.Param("id")
	result := t.DB.Where("id = ? AND user_id = ?", taskID, userID).First(&task)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		}
		return
	}

	task.Status = true

	t.DB.Save(&task)

	c.JSON(http.StatusOK, gin.H{
		"message": "task marked as done",
		"data":    task,
	})
}
