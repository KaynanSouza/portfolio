package routes

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
	controller "trabalhoCaio/controllers"
	"trabalhoCaio/dataBase"
)

func HandlerRequests() {

	dataBase.InitDataBase()

	router := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	//router.Use(cors.Default())

	// Users
	router.POST("/user/login", controller.Login)
	router.GET("/user/:username", controller.UserPage)
	router.POST("/user/signup", controller.AddUser)
	router.PUT("/user/:username/addBooks", controller.AddBooksToUser)
	router.PUT("/user/:username/removeBooks", controller.RemoveBooksFromUser)
	router.PUT("/user/:username", controller.UpdateUser)

	// Books
	router.GET("/books", controller.AllBooks)
	router.GET("/books/:title", controller.OneBook)
	router.POST("/admin/book/addBooks", controller.AddBooks)

	// Admin
	router.POST("/admin/singnup", controller.AddAdmin)
	router.GET("/admin/users", controller.GetUsers)
	router.DELETE("/admin/user/:username", controller.DeleteUser)
	router.PUT("/admin/book/:title", controller.UpdateBook)
	router.DELETE("/admin/book/:title", controller.DeleteBook)

	router.Run("10.233.57.152:8080")
}
