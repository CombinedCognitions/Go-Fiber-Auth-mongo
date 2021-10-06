package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func Db() *mongo.Client {
	err := godotenv.Load()
	if err != nil {
		log.Panicln(err)
	}
	uri := GetURL()
	fmt.Println(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}
	// fmt.Println("Successfully connected and pinged.")
	return client
}

func GetURL() string {
	port, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Println("error on load db port from env:", err.Error())
		port = 27017
	}

	return fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASS"),
		os.Getenv("DATABASE_HOST"),
		port,
		os.Getenv("DATABASE_NAME"))
}
