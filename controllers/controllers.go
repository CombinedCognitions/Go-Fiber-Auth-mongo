package controllers

import (
	"context"
	"fmt"
	"go-fiber-auth/db"
	"go-fiber-auth/models"
	"go-fiber-auth/utix"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userClient = db.Db()
var userCollection = userClient.Database("mydb").Collection("users")

func Save(user *models.User) error { //   save to db

	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ADDED NEW USER ", user.Email)
	return err
}

func GetByEmail(email string) (models.User, error) { // get by email
	var result models.User
	//var userlogin models.User

	err := userCollection.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&result)
	if err != nil {
		fmt.Println(err)
	}

	// fmt.Println(result, " this is the result bb")
	return result, err
}

func GetByKey(key string, value string) (models.User, error) {

	filter := bson.D{{key, value}}
	var res models.User

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}

func GetAll() []models.User {
	cursor, err := userCollection.Find(context.Background(), bson.D{})
	if err != nil {
		log.Panic(err)
	}
	var docs []models.User

	for cursor.Next(context.Background()) {
		var single models.User
		err := cursor.Decode(&single)
		if err != nil {
			log.Panic(err)
		}
		docs = append(docs, single)
	}

	return docs

}

func Delete(id string) (*mongo.DeleteResult, error) {

	_id, err1 := primitive.ObjectIDFromHex(id)
	if err1 != nil {
		panic(err1)
	}

	opts := options.Delete().SetCollation(&options.Collation{})

	res, err := userCollection.DeleteOne(context.Background(), bson.D{{"_id", _id}}, opts)
	if err != nil {
		log.Panic(err)
	}

	return res, err
}

func Update(key string, value string, user models.User) {
	filter := bson.D{{key, value}}

	update := bson.D{{"$set", bson.D{{"password", user.Password}}}}

	_, e := userCollection.UpdateOne(context.Background(), filter, update)
	utix.CheckErorr(e)
	fmt.Println("update sucesss")

}

func Close() error {
	err := userClient.Disconnect(context.Background())
	fmt.Println("db closed")
	utix.CheckErorr(err)
	return err
}

func GetByID(key string, value string) (models.User, error) {
	_id, err1 := primitive.ObjectIDFromHex(value)
	utix.CheckErorr(err1)
	filter := bson.D{{key, _id}}
	var res models.User

	err := userCollection.FindOne(context.Background(), filter).Decode(&res)

	return res, err

}
