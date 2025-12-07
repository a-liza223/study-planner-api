package handlers

import (
	"net/http"
	"study-tracker-api/models"
	"study-tracker-api/storage"

	"github.com/gin-gonic/gin"
)

func GetLessons(c *gin.Context) {
	lessons := storage.GetLessons()
	c.JSON(http.StatusOK, lessons)
}

func GetLessonByID(c *gin.Context) {
	id := c.Param("id")
	lesson, found := storage.GetLessonByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Занятие не найдено"})
		return
	}
	c.JSON(http.StatusOK, lesson)
}

func CreateLesson(c *gin.Context) {
	var newLesson models.Lesson
	if err := c.ShouldBindJSON(&newLesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	created, err := storage.CreateLesson(newLesson)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func UpdateLesson(c *gin.Context) {
	id := c.Param("id")
	var updatedLesson models.Lesson
	if err := c.ShouldBindJSON(&updatedLesson); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	updated, err := storage.UpdateLesson(id, updatedLesson)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteLesson(c *gin.Context) {
	id := c.Param("id")

	err := storage.DeleteLesson(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
