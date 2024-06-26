package routes

import(
	controller "github.com/Ridwan-Al-Mahmud/Let-s-Go/JWT-AUTH/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine){
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}