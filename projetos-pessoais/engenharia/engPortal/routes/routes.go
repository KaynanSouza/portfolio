package routes

import (
	"engPortal/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func HandlerRequests() {

	router := gin.Default()

	//config := cors.Config{
	//	AllowOrigins:     []string{"http://localhost:8080", "http://localhost:5173"},
	//	AllowMethods:     []string{"GET"},
	//	AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
	//	ExposeHeaders:    []string{"Content-Length"},
	//	AllowCredentials: true,
	//	MaxAge:           12 * time.Hour,
	//}

	//router.Use(cors.New(config))

	router.Use(cors.Default())

	api := router.Group("/api")
	{
		api.GET("/projects", controller.GetProjects)
		api.GET("/kanban", controller.GetKanbanData)
		api.GET("/docs", controller.GetDocs)
		api.GET("/download", controller.DownloadFile)

		api.GET("/trainings", controller.GetTrainings)
		api.POST("/trainings", controller.CreateTraining)
		api.PUT("/trainings/:id", controller.UpdateTraining)
		api.DELETE("/trainings/:id", controller.DeleteTraining)
	}

	router.Run("0.0.0.0:8081")
}
