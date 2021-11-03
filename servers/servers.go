package servers

import (
	"GoLang-API-2/handlers"
	"github.com/gin-gonic/gin"
)

func Init(port string) error {
	//create a new gin router
	router := gin.Default()

	//define a single endpoint
	router.GET("/", handlers.HelloListAdmin)

	// CRUD endpoints for data

	//create
	router.POST("/createList", handlers.CreateListItem)

	router.POST("/login", handlers.LoginUser)

	router.POST("/signup", handlers.SignupUser)

	//retrieve
	router.GET("/getListItem/:title", handlers.GetSingleListItem)

	router.GET("/getLists", handlers.GetMultipleListItem)

	router.GET("/getUser", handlers.GetAllUsers)

	//update
	router.PATCH("/updateListItem/:title", handlers.UpdateList)

	// delete
	router.DELETE("/deleteList/:title", handlers.DeleteList)

	err := router.Run(":" + port)
	if err != nil {
		return err
	}
	return nil
}
