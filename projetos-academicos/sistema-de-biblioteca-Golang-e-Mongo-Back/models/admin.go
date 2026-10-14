package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Admin struct {
	ID       primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name     string             `bson:"fullName" json:"fullName"`
	Email    string             `bson:"email" json:"email"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
}

type CreateAdminRequest struct {
	Name     string `bson:"name" json:"name"`
	Email    string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}
