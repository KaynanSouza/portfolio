package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID          primitive.ObjectID `Bson:"_id" json:"id"`
	Title       string             `Bson:"title" json:"title" validate:"nonzero"`
	Author      string             `Bson:"author" json:"author" validate:"nonzero regexp=^[a-zA-Z]*$"`
	Year        int                `Bson:"year" json:"year" validate:"nonzero regexp=^[0-9]*$"`
	Pages       int                `Bson:"pages" json:"pages"`
	Description string             `Bson:"description" json:"description" validate:"nonzero"`
}

type CreateBookRequest struct {
	Title       string `Bson:"title" json:"title"`
	Author      string `Bson:"author" json:"author"`
	Year        int    `Bson:"year" json:"year"`
	Pages       int    `Bson:"pages" json:"pages"`
	Description string `Bson:"description" json:"description"`
}
