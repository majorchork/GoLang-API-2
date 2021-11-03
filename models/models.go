package models

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)
type User struct {
	ID       string    `json:"id" bson:"id"`
	Name     string    `json:"name" bson:"name"`
	Email    string    `json:"email" bson:"email"`
	Password string    `json:"-,omitempty" bson:"password"`
	Ts       time.Time `json:"timestamp" bson:"timestamp"`
}
type List struct {
	ID       string    `json:"id" bson:"id"`
	Title    string    `json:"title" bson:"title"`
	Activity string    `json:"activity" bson:"activity"`
	Executor string    `json:"executor" bson:"executor"`
	Ts       time.Time `json:"timestamp" bson:"timestamp"`
}
type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}