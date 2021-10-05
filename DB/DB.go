package DB

import (
	"GoLang-API-2/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

const (
	DbName         = "listsdb"
	ListCollection = "lists"
	UserCollection = "Users"

	jwtSecret = "secretname"
)
var dbClient *mongo.Client
func init (){
	// establishing connection parameters
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// establishing connection detail and pointing to the db
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil{
		log.Fatalf("could not connect to the db: %v/n", err)
	}
	dbClient = client
	// after connection to the db, this checks if the db is active and available at the moment
	err = dbClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Mongo db not available: %v/n", err)
	}
}
func CreateList(list *models.List)(*models.List, error){
	_, err := dbClient.Database(DbName).Collection(ListCollection).InsertOne(context.Background(), list)

	return list, err
}
// reaserch about & in handlers
func GetSingleListItem(title string)(*models.List, error){
	var list models.List
	query := bson.M{
		"title": title,
	}
	err := dbClient.Database(DbName).Collection(ListCollection).FindOne(context.Background(), query).Decode(&list)
	if err != nil{
		return nil, err
	}
	return &list, nil
}
func GetAllUserHandler() ([]models.User, error) {
	var users []models.User
	cursor, err := dbClient.Database(DbName).Collection(UserCollection).Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func GetMultipleListItem() ([]models.List, error){
	var lists []models.List
	cursor, err := dbClient.Database(DbName).Collection(ListCollection).Find(context.Background(), bson.M{})
	if err != nil{
		return nil, err
	}
	err = cursor.All(context.Background(), &lists)
	if err != nil {
		return nil, err
	}
	return lists, nil
}
func UpdateListItem(title, activity, executor string)(*models.List, error){

	filterQuery := bson.M{
		"title": title,
	}
	updateQuery := bson.M{
		"$set": bson.M{
			"title": title,
			"activity": activity,
			"executor": executor,
		},
	}
	_, err := dbClient.Database(DbName).Collection(ListCollection).UpdateOne(context.Background(), filterQuery, updateQuery)
	if err != nil {
		return nil, err
	}
	return nil, err
}
func DeleteListItem(title string) error{
	query := bson.M{
		"title" : title,
		//"executor": executor,
	}
	_, err := dbClient.Database(DbName).Collection(ListCollection).DeleteOne(context.Background(),query)
	if err != nil {
		return err
	}
	return nil
}



func LoginHandler(Email string) (*models.User, error) {
	var user models.User

	query := bson.M{
		"email": Email,
	}
	err := dbClient.Database(DbName).Collection(UserCollection).FindOne(context.Background(), query).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, err

}


func SignUpCheckUserExists (email string) bool {
	query := bson.M{
		"email": email,
	}
	check, err := dbClient.Database(DbName).Collection(ListCollection).CountDocuments(context.Background(), query)
		if err != nil {
			return false
		}
	// if the check is greater than zero that means a user exists already with that email
	if check > 0 {

			return true
	}
	return false

}
func CreateUser(user *models.User)(*models.User, error){
	_, err := dbClient.Database(DbName).Collection(UserCollection).InsertOne(context.Background(), user)

	return user, err
}
