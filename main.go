package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"strings"
	"time"
)

const(
	DbName = "listsdb"
	ListCollection = "lists"
	UserCollection = "Users"

	jwtSecret = "secretname"
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


var dbClient *mongo.Client
func main (){
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

	_ = router.GET("/getUsers", getAllUserHandler)

	// update
	_ = router.PATCH("/updateListItem/:title", updateListItem)

	// delete
	_ = router.DELETE("/deleteListItem/:title", deleteListItem)

	// post
	_ = router.POST("/login", loginHandler)

	_ = router.POST("/signup", signupHandler)

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
	//  we need to find out the identity of the user
	//  as this endpoint does not request for the users details like email or user id
	// as you've already gotten the details during login
	//  the request only contains the task and the jwt token

	//  the jwt token is what we use to identify the user
	// we generate this token during login or signup
	// because it is at that point that we confirm things like password and other security details we might be intersted in
	// you can't be asking the user for password at every endpoint
	// the jwt only contains the things we put inside
	// the only thing we need for our app to identify the user is the users id

	// for http request, the standard way the jwt is usually sent is as a request header
	// we need to get jwt token from request header using then key
	// for the jwt the key name is "Authorization"
	// get Jwt token from request
	authorization := c.Request.Header.Get("Authorization")
	fmt.Println(authorization)
	
	// we return an error to the user if the token was not supplied
	if authorization == ""{
		c.JSON(401, gin.H{
			"error": "authentication token not supplied",
		})
		return
	}


	jwtToken := ""
	// split the authenthication token which looks like "Bearer asdsadsdsdsdsa........."
	//  so that we can get the second part of the string which is the actual jwt token
	splitTokenArray := strings.Split(authorization, " ")
	if len(splitTokenArray) > 1 {
		jwtToken = splitTokenArray[1]
	}

	// create an empty claims object to store the claims(userid,.......)
	// decode token to get claims
	claims := &Claims{}

	keyFunc := func(token *jwt.Token) (i interface{}, e error) {
		return []byte(jwtSecret), nil
	}

	token, err := jwt.ParseWithClaims(jwtToken, claims, keyFunc)
	if !token.Valid {
		c.JSON(400, gin.H{
			"error": "invalid jwt token",
		})
		return
	}
	//create a list item
	var list List

	// gets the list item from client (postman)
	// fills up our empty list item with sent data
	err = c.ShouldBindJSON(&list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}
	// generate task id
	listId := uuid.NewV4().String()

	list = List{
		ID:       listId,
		Title:    list.Title,
		Activity: list.Activity,
		Executor: claims.UserId,
		Ts:       time.Now(),
	}


	// add single item to list of Items
	// Lists = append(Lists, list)
	// linking to a db
	_, err = dbClient.Database(DbName).Collection(ListCollection).InsertOne(context.Background(),list)
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
func getAllUserHandler(c *gin.Context){
	var Users []User
	cursor, err := dbClient.Database(DbName).Collection(UserCollection).Find(context.Background(), bson.M{})
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, user not found",
		})
		return
	}
	err = cursor.All(context.Background(), &Users)
	if err != nil {
		c.JSON(500, gin.H{
			"error" : "unable to process request, user not found",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "list item found",
		"data": Users,
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