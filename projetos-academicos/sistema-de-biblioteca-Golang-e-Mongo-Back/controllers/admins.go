package controller

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"net/http"
	"time"
	"trabalhoCaio/dataBase"
	"trabalhoCaio/models"
)

func AddAdmin(c *gin.Context) {
	var body models.CreateAdminRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Erro": "Invalid request body",
		})
		return
	}

	usersCollection := dataBase.DB.Collection("Admins")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verificar se o usuário já existe (por nome de usuário ou e-mail)
	var existingAdmin models.Admin

	err := usersCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"username": body.Username},
			{"email": body.Email},
		},
	}).Decode(&existingAdmin)

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
	admin := models.Admin{
		ID:       primitive.NewObjectID(),
		Name:     body.Name,
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
	}

	_, err = usersCollection.InsertOne(ctx, admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Erro": "Unable to add user",
		})
		return
	}

	c.JSON(http.StatusCreated, admin)
}
