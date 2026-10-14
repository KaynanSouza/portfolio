package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
	"time"
	"trabalhoCaio/dataBase"
	"trabalhoCaio/models"
)

// Função para deletar um livro pelo título.
func DeleteBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booksCollection := dataBase.DB.Collection("Books")

	title := c.Params.ByName("title")

	titleConv := strings.Replace(title, "%20", " ", -2)

	filter := bson.D{{"title", titleConv}}

	result, err := booksCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete book"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No book was found with the given title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}

// Função para atualizar um livro
func UpdateBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booksCollection := dataBase.DB.Collection("Books")

	title := c.Params.ByName("title")
	var updatedBook models.Book

	if err := c.BindJSON(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	titleConv := strings.Replace(title, "%20", " ", -2)

	filter := bson.D{{"title", titleConv}}
	update := bson.D{
		{"$set", bson.D{
			{"author", updatedBook.Author},
			{"year", updatedBook.Year},
			{"pages", updatedBook.Pages},
			{"description", updatedBook.Description},
		}},
	}

	result, err := booksCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update book"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No book was found with the given title"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book updated successfully"})
}

// Mostrar somente os dados de um livro
func OneBook(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booksCollection := dataBase.DB.Collection("Books")

	title := c.Params.ByName("title")
	var book bson.M

	titleConv := strings.Replace(title, "%20", " ", -2)

	err := booksCollection.FindOne(ctx, bson.D{{"title", titleConv}}).Decode(&book)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No book found with the specified title",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occurred while fetching the book",
			})
		}
		return
	}

	c.JSON(http.StatusOK, book)
}

// Função para mostrar todos os livros
func AllBooks(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	booksCollection := dataBase.DB.Collection("Books")

	cursor, err := booksCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch users",
		})
		return
	}
	defer cursor.Close(ctx)

	var books []models.Book
	for cursor.Next(ctx) {
		var book models.Book
		if err := cursor.Decode(&book); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to decode user",
			})
			return
		}
		books = append(books, book)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cursor error",
		})
		return
	}

	c.JSON(http.StatusOK, books)
}

// Entrada de livros no banco de dados
func AddBooks(c *gin.Context) {
	var body models.CreateBookRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Erro": "Invalid request body",
		})
		return
	}

	// Conectar com o banco
	booksCollection := dataBase.DB.Collection("Books")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar se o livro ja existe
	var existingBook models.Book
	err := booksCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"title": body.Title},
		},
	}).Decode(&existingBook)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"Erro": "Book already exists",
		})
		return
	}

	book := models.Book{
		Title:       body.Title,
		Author:      body.Author,
		Year:        body.Year,
		Pages:       body.Pages,
		Description: body.Description,
	}

	res, err := booksCollection.InsertOne(ctx, book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Unable to add user",
		})
		return
	}

	book.ID = res.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, book)
}
