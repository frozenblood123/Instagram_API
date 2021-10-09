package post_user

import (
	"context"
	"insta/user_data"
	"insta/user_validation"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PostUserHandler struct {
	postCollection *mongo.Collection
}

func NewPostUserHandler(col *mongo.Collection) *PostUserHandler {
	return &PostUserHandler{
		postCollection: col,
	}
}

func (h *PostUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			//id from url
			userId := r.URL.Path[len("/posts/users/"):]

			postCursor, err := h.postCollection.Find(context.Background(), bson.D{{"userId", userId}})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			posts := &[]user_data.OutPost{}
			err = postCursor.All(context.Background(), posts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			user_validation.WriteJson(w, r, posts)
		}

	default:
		{
			w.Write([]byte("Method not implemented"))
		}
	}

}
