package routes

import(
	controller "github.com/Ridwan-Al-Mahmud/Let-s-Go/JWT-AUTH/controllers"
	"github.com/Ridwan-Al-Mahmud/Let-s-Go/JWT-AUTH/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}