package handlers

import (
	"net/http"
	"study-tracker-api/models"
	"study-tracker-api/storage"

	"github.com/gin-gonic/gin"
)

func GetAssignments(c *gin.Context) {
	assignments := storage.GetAssignments()
	c.JSON(http.StatusOK, assignments)
}

func GetAssignmentByID(c *gin.Context) {
	id := c.Param("id")
	assignment, found := storage.GetAssignmentByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "Задание не найдено"})
		return
	}
	c.JSON(http.StatusOK, assignment)
}

func CreateAssignment(c *gin.Context) {
	var newAssignment models.Assignment
	if err := c.ShouldBindJSON(&newAssignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	created, err := storage.CreateAssignment(newAssignment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	var updatedAssignment models.Assignment
	if err := c.ShouldBindJSON(&updatedAssignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	updated, err := storage.UpdateAssignment(id, updatedAssignment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func DeleteAssignment(c *gin.Context) {
	id := c.Param("id")

	err := storage.DeleteAssignment(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
