package router

import (
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/controllers"
	"task-5-vix-btpns-HaiqalRamanizarAlFajri/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.POST("/users/register", controllers.Register)
	r.GET("/users/login", controllers.Login)
	r.PUT("/users/:userId", controllers.UpdateUser)
	r.DELETE("/users/:userId", controllers.DeleteUser)

	r.GET("/photos", controllers.GetPhotos)

	secured := r.Group("/").Use(middlewares.Auth())
	{
		secured.POST("/photos", controllers.CreatePhoto)
		secured.PUT("/photos/:photoId", controllers.UpdatePhoto)
		secured.DELETE("/photos/:photoId", controllers.DeletePhoto)
	}

	return r
}
