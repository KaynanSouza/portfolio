package dataBase

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"time"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
)

func InitDataBase() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}

	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Erro ao fazer ping no MongoDB: %v", err)
	}

	Client = client
	DB = client.Database("Biblioteca")

	log.Println("Conexão ao MongoDB estabelecida com sucesso.")
}
