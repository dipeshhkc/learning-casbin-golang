package route

import (
	"casbin-golang/controller"
	"casbin-golang/middleware"
	"casbin-golang/repository"
	"fmt"
	"log"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//SetupRoutes : all the routes are defined here
func SetupRoutes(db *gorm.DB) {
	httpRouter := gin.Default()

	// Initialize  casbin adapter
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize casbin adapter: %v", err))
	}

	userRepository := repository.NewUserRepository(db)

	if err := userRepository.Migrate(); err != nil {
		log.Fatal("User migrate err", err)
	}

	userController := controller.NewUserController(userRepository)

	apiRoutes := httpRouter.Group("/api")

	{
		apiRoutes.POST("/register", userController.AddUser)
		apiRoutes.POST("/signin", userController.SignInUser)
	}

	userProtectedRoutes := apiRoutes.Group("/users", middleware.AuthorizeJWT())
	{
		userProtectedRoutes.GET("/", middleware.Authorize("report", "read", adapter), userController.GetAllUser)
		userProtectedRoutes.GET("/:user", middleware.Authorize("report", "read", adapter), userController.GetUser)
		userProtectedRoutes.PUT("/:user", middleware.Authorize("report", "write", adapter), userController.UpdateUser)
		userProtectedRoutes.DELETE("/:user", middleware.Authorize("report", "write", adapter), userController.DeleteUser)
	}

	httpRouter.Run()

}
