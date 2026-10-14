package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `bson:"fullName" json:"fullName"`
	Email    string             `bson:"email" json:"email"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Books    []string           `bson:"books" json:"books"`
}

type CreateUserRequest struct {
	Name     string `bson:"fullName" json:"fullName"`
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type LoginRequest struct {
	User     string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
