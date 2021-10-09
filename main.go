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
	// Initializing the mongodb client
	mongoClient := connectToDb()

	// connecting to collections
	userCollection := mongoClient.Database("Insta").Collection("Users")
	postCollection := mongoClient.Database("Insta").Collection("Posts")

	//Initializing habdlers from each route
	userHandler := post_user.NewUserHandler(userCollection)
	postHandler := post_user.NewPostHandler(postCollection)
	postUserHandler := post_user.NewPostUserHandler(postCollection)

	//Routing the server using different handlers
	http.Handle("/users/", userHandler)
	http.Handle("/posts/", postHandler)
	http.Handle("/posts/users/", postUserHandler)

	// Starting the server at localhost:8080
	fmt.Println("Starting the server at localhost:8080 ...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

// Function to connect to the mongodb server
func connectToDb() *mongo.Client {
	// Connection to db at localhost:27017
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
