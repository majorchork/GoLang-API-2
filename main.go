package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type User struct {
	ID string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Activity string `json:"activity" bson:"activity"`
	Executor string `json:"executor" bson:"executor"`
	Ts time.Time `json:"timestamp" bson:"timestamp"`
}
type List struct {
	ID string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title"`
	Activity string `json:"activity" bson:"activity"`
	Executor string `json:"executor" bson:"executor"`
	Ts time.Time `json:"timestamp" bson:"timestamp"`
}
type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}
	// creating an empty array of list
var Lists []List


var db
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
	_ = router.GET("/getListItem/:title", getSingleListItem)

	_ = router.GET("/getListItems", getMultipleListItem)

	// update
	_ = router.PATCH("/updateListItem/:title", updateListItem)

	// delete
	_ = router.DELETE("/deleteListItem/:title", deleteListItem)

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
	// linking to a db
	_, err = dbClient.Database("listsdb").Collection("lists").InsertOne(context.Background(),list)
	if err != nil{
		fmt.Println("error creating list", err)
		// if saving failed
		c.JSON(500, gin.H{
			"error": "could not create list, unable to process request",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "succesfully created list item",
		"data":    list,
	})
	}
func getSingleListItem(c *gin.Context) {
	title := c.Param("title")

	fmt.Println("title", title)

	var list List
	/* listAvailable := false
	for _, value := range Lists{
		//check the current iteration of list items
		// check for a match with client request
		if value.Title == title {
			// if it matches the asign the value to the empty list item we created and display
			list  = value
		listAvailable = true

		}
	}
	// if no match was found
	// the empty list we created would still be empty
	// check if the user is empty if so return a not found error
	if !listAvailable {
		c.JSON(404, gin.H{
			"error": "no list with tittle found:" + title,
		})
		return
	}*/
	// linking to a db
	query := bson.M{
		"title" : title,
	}
	// ask vic about _, why it is useful up and not now and
	err := dbClient.Database("listsdb").Collection("lists").FindOne(context.Background(),query).Decode(&list)
	if err != nil{
		fmt.Println("list not found", err)
		// if saving failed
		c.JSON(400, gin.H{
			"error": "could not find list" + title,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "list item found",
		"data": list,
	})
}
func getMultipleListItem(c *gin.Context){
	var lists []List
	cursor, err := dbClient.Database("listsdb").Collection("lists").Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, list not found",
		})
		return
	}
	err = cursor.All(context.Background(), &lists)
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, list not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "list item found",
		"data": lists,
	})
}
func updateListItem(c *gin.Context){
	title := c.Param("title")

	var list List
	err := c.ShouldBindJSON(&list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}

	filterQuery := bson.M{
		"title": title,
	}

	updateQuery := bson.M{
		"$set": bson.M{
			"title": list.Title,
			"activity": list.Activity,
			"executor": list.Executor,

		},
	}

	_, err = dbClient.Database("listsdb").Collection("lists").UpdateOne(context.Background(), filterQuery, updateQuery)
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, Update failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Update Successful",

	})

}
func deleteListItem(c *gin.Context){
	title := c.Param("title")

	Query := bson.M{
		"title": title,
	}


	_, err := dbClient.Database("listsdb").Collection("lists").DeleteOne(context.Background(), Query)
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, delete failed",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "Successfully Deleted",
	})

}