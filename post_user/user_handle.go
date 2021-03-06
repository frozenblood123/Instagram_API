package post_user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"insta/user_data"
	"insta/user_validation"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userCollection *mongo.Collection
}

func NewUserHandler(col *mongo.Collection) *UserHandler {
	return &UserHandler{
		userCollection: col,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			h.getUser(w, r)
		}
	case http.MethodPost:
		{
			h.createUser(w, r)
		}
	default:
		{
			http.Error(w, "Method not implemented", http.StatusMethodNotAllowed)
		}
	}
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	// Getting the request body into the user data
	user := &user_data.InUser{}
	ok := user_validation.ReadJson(w, r, user)
	if !ok {
		return
	}

	//validating user credential
	err1 := user_validation.ValidateUser(user)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	// Hashing password
	hashedPassword := sha256.New()
	hashedPassword.Write([]byte(user.Password))
	user.Password = fmt.Sprintf("%x\n", hashedPassword.Sum(nil))

	//Inserting user into mongodb
	userResult, err := h.userCollection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Successfully created user
	w.Write([]byte(fmt.Sprintf("Successfully created user with id: %v", userResult.InsertedID)))
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	// Getting id from url
	id := r.URL.Path[len("/users/"):]
	fmt.Println(id)

	// Retrieve user data from db and store in user
	user := &user_data.OutUser{}
	userResult := h.userCollection.FindOne(context.Background(), bson.D{{"_id", id}})
	err := userResult.Decode(user)
	if err != nil {
		w.Write([]byte("unable to get data"))
	} else {
		// Marshal user data to json and send to client
		user_validation.WriteJson(w, r, user)
	}

}
