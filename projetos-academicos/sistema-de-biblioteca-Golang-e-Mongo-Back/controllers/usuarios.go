package controller

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
	"trabalhoCaio/dataBase"
	"trabalhoCaio/models"
)

// Função para buscar somente um usuario
func UserPage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	username, err := c.Cookie("username")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not logged in",
		})
		return
	}

	var user bson.M

	usernameConv := strings.Replace(username, "%20", " ", -2)

	erro := usersCollection.FindOne(ctx, bson.D{{"username", usernameConv}}).Decode(&user)
	if erro != nil {
		if erro == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No user found with the specified username",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error occurred while fetching the user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// Função para buscar todos os usuarios
func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	cursor, err := usersCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to fetch users",
		})
		return
	}
	defer cursor.Close(ctx)

	var users []models.User
	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to decode user",
			})
			return
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Cursor error",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

// Função para adicionar um usuario
func AddUser(c *gin.Context) {
	var body models.CreateUserRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Erro": "Invalid request body",
		})
		return
	}

	usersCollection := dataBase.DB.Collection("Users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar se o usuário já existe (por nome de usuário ou e-mail)
	var existingUser models.User

	err := usersCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": body.Username},
			{"email": body.Email},
		},
	}).Decode(&existingUser)

	if err == nil {
		// Encontrou um usuário existente
		c.JSON(http.StatusConflict, gin.H{
			"Erro": "User already exists",
		})
		return
	} else if err != mongo.ErrNoDocuments {
		// Ocorreu um erro que não é "documento não encontrado"
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Error checking for existing user",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Error generating password hash",
		})
		return
	}

	// Adicionar novo usuário se não existir
	user := models.User{
		ID:       primitive.NewObjectID(),
		Name:     body.Name,
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
		Books:    []string{},
	}

	_, err = usersCollection.InsertOne(ctx, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Unable to add user",
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Função para autenticar usuário e admin
//func Login(c *gin.Context) { // Função para lidar com o login
//
//	var login models.LoginRequest
//	if err := c.ShouldBindJSON(&login); err != nil { // Faz o bind do JSON de entrada para infoLogin
//		log.Printf("Erro ao fazer bind do JSON admin: %v", err)
//		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
//		return
//	}
//
//	LoginUsuario(c, login)
//	LoginAdmin(c, login)
//}

// Função para update do usuario
func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	username := c.Params.ByName("username")
	var updatedUser models.CreateUserRequest

	if err := c.BindJSON(&updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameConv := strings.Replace(username, "%20", " ", -2)

	filter := bson.D{{"username", usernameConv}}
	update := bson.D{
		{"$set", bson.D{
			{"name", updatedUser.Name},
			{"email", updatedUser.Email},
			{"password", updatedUser.Password},
		}},
	}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user was found with the given username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// Função para deletar usuario
func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	username := c.Params.ByName("username")

	usernameConv := strings.Replace(username, "%20", " ", -2)

	filter := bson.D{{"username", usernameConv}}

	result, err := usersCollection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to delete user"})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user was found with the given username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// Função para adicionar livros ao usuario
func AddBooksToUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	username := c.Params.ByName("username")
	var books []string

	if err := c.BindJSON(&books); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameConv := strings.Replace(username, "%20", " ", -2)

	filter := bson.D{{"username", usernameConv}}
	update := bson.D{
		{"$addToSet", bson.D{
			{"books", bson.D{{"$each", books}}},
		}},
	}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add books to user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user was found with the given username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Books added successfully"})
}

// Função para remover livros do usuario
func RemoveBooksFromUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := dataBase.DB.Collection("Users")

	username := c.Params.ByName("username")
	var books []string

	if err := c.BindJSON(&books); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usernameConv := strings.Replace(username, "%20", " ", -2)

	filter := bson.D{{"username", usernameConv}}
	update := bson.D{
		{"$pull", bson.D{
			{"books", bson.D{{"$in", books}}},
		}},
	}

	result, err := usersCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to remove books from user"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No user was found with the given username"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Books removed successfully"})
}

// Função para autenticar usuário
func Login(c *gin.Context) {
	var loginInfo models.LoginRequest
	if err := c.ShouldBindJSON(&loginInfo); err != nil {
		log.Printf("Erro ao fazer bind do JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usersCollection := dataBase.DB.Collection("Users")

	filter := bson.D{{"username", loginInfo.User}}
	var result models.User

	// Verificar se o usuário existe pelo nome de usuário
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := usersCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("Usuário não encontrado: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
			return
		}
		log.Printf("Erro ao encontrar usuário: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while finding user"})
		return
	}

	// Compare the password with the hashed password stored in the database
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(loginInfo.Password))
	if err != nil {
		log.Printf("Erro ao comparar senha: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid login credentials"})
		return
	}

	// Definir o cookie do usuário
	//c.SetCookie("username", loginInfo.User, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    result,
	})
}
