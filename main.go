package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type List struct {
	Title string `json:"title"`
	Activity string `json:"activity"`
	Executor string `json:"executor"`
}
	// creating an empty array of list
var Lists []List


var dbClient *mongo.Client

func main() {

	// creating a mongo.Client then connect function under
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	//we are trying to connect to mongodb on a specified url
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("could not conect to the db: %v\n")
	}
	dbClient = client
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("could not conect to the db: %v\n", err)
	}
	 router := gin.Default()
	// define a single endpoint

	router.GET("/", helloListAdmin)
	// crude endpoints for data

	// create
	_ = router.POST("/createListItem", createListItem)

	// retrieve
	_ = router.GET("/getListItem/:tittle", getSingleListItem)

	//_ = router.GET("/getListItems", getmultipleListItem)

	// update
	//_ = router.PATCH("/updateListItem/:name", updateListItem)

	// delete
	//_ = router.DELETE("/deleteListItem/:name", deleteListItem)

	// run the server on the port 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	_ = router.Run(":" + port)

}
func helloListAdmin(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Welcome Please create your to do list",
	})
}
func createListItem(c *gin.Context)  {
	//create a list item
	var list List

	// gets the list item from client (postman)
	// fills up our empty list item with sent data
	err := c.ShouldBindJSON(&list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}
	// add single item to list of Items
	Lists = append(Lists, list)

	c.JSON(200, gin.H{
		"message": "succesfully created list item",
		"data":    list,
	})
	}
func getSingleListItem(c *gin.Context) {
	tittle := c.Param("tittle")

	fmt.Println("tittle", tittle)

	var list List
	for _, value := range Lists{
		//check the current iteration of list items
		// check for a match with client request
		if value.Title == tittle {
			// if it matches the aasign the value to the empty list item we created and display
			list  = value

		}
	}
	// if no match was found
	// the empty list we creaated would still be empty
	// check if the user is empty if so return a not found error
	if &list == nil {
		c.JSON(404, gin.H{
			"error": "no user with name found:" + tittle,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": Lists,
	})
}