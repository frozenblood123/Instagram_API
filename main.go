package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"insta/post_user"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//connecting to client
	mongoClient := connectToDb()

	//connecting to collections
	userCollection := mongoClient.Database("Insta").Collection("Users")
	postCollection := mongoClient.Database("Insta").Collection("Posts")

	userHandler := post_user.NewUserHandler(userCollection)
	postHandler := post_user.NewPostHandler(postCollection)
	postUserHandler := post_user.NewPostUserHandler(postCollection)

	http.Handle("/users/", userHandler)
	http.Handle("/posts/", postHandler)
	http.Handle("/posts/users/", postUserHandler)

	fmt.Println("Starting the server at localhost:8080 ...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

//connecting to mongodb server
func connectToDb() *mongo.Client {
	mongoClient, err := mongo.Connect(context.Background(), &options.ClientOptions{
		Auth: &options.Credential{
			Username: "mongodb",
			Password: "asdfghjkl",
		},
	})
	if err != nil {
		log.Fatalf("Unable to connect to db\n[Error]: %v", err)
	}

	//Creating user and post collections
	mongoClient.Database("Insta").CreateCollection(context.Background(), "Users")
	mongoClient.Database("Insta").CreateCollection(context.Background(), "Posts")

	return mongoClient
}
