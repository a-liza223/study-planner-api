package main

import (
	"study-tracker-api/handlers"
	"study-tracker-api/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализируем хранилище
	storage.InitStorage()

	router := gin.Default()

	// Группа маршрутов для занятий
	lessonGroup := router.Group("/api/v1/lessons")
	{
		lessonGroup.GET("", handlers.GetLessons)
		lessonGroup.GET("/:id", handlers.GetLessonByID)
		lessonGroup.POST("", handlers.CreateLesson)
		lessonGroup.PUT("/:id", handlers.UpdateLesson)
		lessonGroup.DELETE("/:id", handlers.DeleteLesson)
	}

	// Группа маршрутов для заданий
	assignmentGroup := router.Group("/api/v1/assignments")
	{
		assignmentGroup.GET("", handlers.GetAssignments)
		assignmentGroup.GET("/:id", handlers.GetAssignmentByID)
		assignmentGroup.POST("", handlers.CreateAssignment)
		assignmentGroup.PUT("/:id", handlers.UpdateAssignment)
		assignmentGroup.DELETE("/:id", handlers.DeleteAssignment)
	}

	// Маршрут для просмотра нагрузки (как в твоём консольном приложении)
	router.GET("/api/v1/workload", handlers.GetWorkload)

	// Запуск сервера
	router.Run(":3000")
}
