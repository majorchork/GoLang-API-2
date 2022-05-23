package handlers


import (
	"GoLang-API-2/DB"
	"GoLang-API-2/jawt"
	"GoLang-API-2/models"
		"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"

	"golang.org/x/crypto/bcrypt"


	"strings"
	"time"

)

const (
 jwtSecret = "secretname"
	)


func HelloListAdmin (c *gin.Context){
	c.JSON(200, gin.H{
		"welcome": "please create your list",
	})
}
func CreateListItem (c *gin.Context){

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


	// we return an error to the user if the token was not supplied
	if authorization == "" {
		c.JSON(401, gin.H{
			"error": "authentication token not supplied",
		})
		return
	}

	jwtToken := ""
	// split the authentication token which looks like "Bearer asdsadsdsdsdsa........."
	//  so that we can get the second part of the string which is the actual jwt token
	splitTokenArray := strings.Split(authorization, "")
	if len (splitTokenArray) > 1 {
		jwtToken = splitTokenArray[1]
	}
	// create an empty claims object to store the claims(userid,.......)
	// decode token to get claims
	claims, err := jawt.ValidateToken(jwtToken)
	if err != nil{
		c.JSON(401, gin.H{
			"error": "invalid authentication token",
		})
		return
	}
	// now that we have validated the token we can continue to create the list

	var list models.List
	// generate task id
	listId := uuid.NewV4().String()

	list = models.List{
		ID: listId,
		Title: list.Title,
		Activity: list.Activity,
		Executor: claims.UserId,
		Ts: time.Now(),
	}
	// add single item to list of Items on the Db
	_, err = DB.CreateList(&list)
	if err != nil {
		fmt.Println("error creating taask", err)
		c.JSON(500, gin.H{
			"error": "could not process request, list not created",
		})
		return
	}
	 c.JSON(200, gin.H{
	 	"message": "Succesfully created list",
	 	"data": list,
	 })

}
func GetSingleListItem(c *gin.Context){
	title := c.Param("title")

	fmt.Println("title", title)
	list, err := DB.GetSingleListItem("title")
	if err != nil {
		fmt.Println("title not found", err)
		c.JSON(404, gin.H{
			"error": "invalid title:" + title,
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"Data": list,
	})
}
func GetAllUsers (c *gin.Context){
	users, err := DB.GetAllUserHandler()
	if err != nil {
		c.JSON(500, gin.H{
			"error":"Unable to Process request, could not get users",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "success",
		"data": users,
	})
}
func GetMultipleListItem (c *gin.Context){
	authorization := c.Request.Header.Get("Authorization")

	if authorization == "" {
		c.JSON(401, gin.H{
			"error": "authentication token not supplied",
		})
		return
	}

	jwtToken := ""
	splitTokenArray := strings.Split(authorization, "")
	if len (splitTokenArray) > 1 {
		jwtToken = splitTokenArray[1]
	}
	claims, err := jawt.ValidateToken(jwtToken)
	if err != nil {
		c.JSON(401, gin.H{
				"error": "invalid auth token supplied",
		})
		return
	}

	lists, err := DB.GetMultipleListItem(claims.UserId)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "could not process request, unable to get lists",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "sucess",
		"data": lists,
	})
}
func UpdateList (c *gin.Context){
	// get the value passed by the client
	title := c.Param("title")

	// create an empty object to store the request data
	var list models.List
	// get the user data sent from the client
	//fills up our empty user object with sent data
	err := c.ShouldBindJSON(&list)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}
	_, err = DB.UpdateListItem(title, list.Executor, list.Activity)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "error processing request, unable to update list",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "list updated successfully!",
	})
}
func DeleteList (c *gin.Context)  {
	authorization := c.Request.Header.Get("authentication")
	if authorization == ""{
		c.JSON(401, gin.H{
			"error": "authentication token not supplied",
		})
		return
	}
	jwtToken := ""
	splitTokenArray := strings.Split(authorization, "")
	if len (splitTokenArray) > 1{
		jwtToken = splitTokenArray[1]
	}
	claims, err := jawt.ValidateToken(jwtToken)
	if err != nil {
		c.JSON(401, gin.H{
			"error": "invalid jwt token",
		})
		return
	}
	title := c.Param("title")

	err = DB.DeleteListItem(title, claims.UserId)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "error processing request, unable to delete list item",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "list item successfully deleted",
	})

}
func LoginUser(c *gin.Context) {
	loginReq := struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		fmt.Printf("error getting user from the db: %v/n", err)
		c.JSON(500, gin.H{
			"error": "error processing request, could not get user",
		})
		return
	}
	user, err := DB.LoginHandler(loginReq.Email)
	if err != nil {
		fmt.Printf("error getting user from db: %v/n", err)
		c.JSON(500, gin.H{
			"error":"could not process request, could not get user",
		})
		return
	}
	// if found compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))
	if err != nil {
		fmt.Printf("error validating password: %v/n", err)
		c.JSON(500, gin.H{
			"error":"Invalid login details",
		})
		return
	}
	jwtTokenString, err := jawt.CreateToken(user.ID)

	c.JSON(200, gin.H{
		"message": "signup successful",
		"token": jwtTokenString,
		"data": user,
	})
}
func SignupUser (c *gin.Context){
	type SignupRequest struct {
		Name string `json:"name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var SignupReq SignupRequest

	err := c.ShouldBindJSON(&SignupReq)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "invalid request data",
		})
		return
	}
	exists := DB.SignUpCheckUserExists(SignupReq.Email)
	if exists {
		c.JSON(500, gin.H{
			"error": "Email already exists, Please login or use a different email",
		})
		return
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(SignupReq.Password), bcrypt.DefaultCost)
	hashPassword := string(bytes)

	// generate user id
	userId := uuid.NewV4().String()

	user := models.User{
		ID:       userId,
		Name:     SignupReq.Name,
		Email:    SignupReq.Email,
		Password: hashPassword,
		Ts:       time.Now(),
	}

	// store the users data
	_, err = DB.CreateUser(&user)
	if err != nil {
		fmt.Println("error saving user", err)
		c.JSON(500, gin.H{
			"error": "could not process request, User not saved",
		})
		return
	}
	 // claims are the data that you want to store inside the jwt token
	//	// so whenever someone gives you a token you can decode it and get back this same claims
	claims := models.Claims{
		UserId:        user.ID ,
		StandardClaims: jwt.StandardClaims{
			IssuedAt: time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1 ).Unix(),
		},
	}

	// generate the jwt token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtTokenString, err :=token.SignedString([]byte(jwtSecret))

	c.JSON(200, gin.H{
		"message": "Signup Successful",
		"token": jwtTokenString,
		"data": user,
	})
}