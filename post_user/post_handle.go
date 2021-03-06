package post_user

import (
	"context"
	"fmt"
	"insta/user_data"
	"insta/user_validation"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostHandler struct {
	postCollection *mongo.Collection
}

func NewPostHandler(col *mongo.Collection) *PostHandler {
	return &PostHandler{
		postCollection: col,
	}
}

func (h *PostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getPost(w, r)
	case http.MethodPost:
		h.createPost(w, r)
	default:
		http.Error(w, "Bad request", http.StatusBadRequest)
	}
}

func (h *PostHandler) createPost(w http.ResponseWriter, r *http.Request) {
	//post from request
	post := &user_data.InPost{}
	ok := user_validation.ReadJson(w, r, post)
	if !ok {
		return
	}

	//validation of post
	if err := user_validation.ValidatePost(post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generating id for post
	rand.Seed(time.Now().UnixNano())
	post.Id = strconv.FormatInt(int64(rand.Uint64()), 10)
	// Generating time stamp
	post.PostedTimestamp = time.Now()

	// Inserting into the db
	_, err := h.postCollection.InsertOne(context.Background(), post)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// Post successfully created
	w.Write([]byte("successfully created post"))
}

func (h *PostHandler) getPost(w http.ResponseWriter, r *http.Request) {
	//id from url
	id := r.URL.Path[len("/posts/"):]
	fmt.Println(id)

	post := &user_data.OutPost{}
	//post from db
	err := h.postCollection.FindOne(context.Background(), bson.D{{"_id", id}}).Decode(post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// showing the post to user
	user_validation.WriteJson(w, r, post)
}
