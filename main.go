package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
func encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

type User struct {
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
}

var client *mongo.Client

type userHandlers struct {
	sync.Mutex
	store map[string]User
}

func (h *userHandlers) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
		return
	case "POST":
		h.post(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}
func (h *userHandlers) get(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	var users []User
	collection := client.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var user User
		cursor.Decode(&user)
		user.Password = string(decrypt([]byte(user.Password), "password"))
		users = append(users, user)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(users)
}
func (h *userHandlers) getUser(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, _ := primitive.ObjectIDFromHex(parts[2])
	var user User
	collection := client.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{Id: id}).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(user)
}
func (h *userHandlers) post(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	collection := client.Database("test").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user.Password = string(encrypt([]byte(user.Password), "password"))
	result, _ := collection.InsertOne(ctx, user)
	json.NewEncoder(w).Encode(result)
}

func newuserHandlers() *userHandlers {
	return &userHandlers{
		store: map[string]User{},
	}
}

type Post struct {
	User     string             `json:"user,omitempty" bson:"user,omitempty"`
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string             `json:"caption,omitempty" bson:"caption,omitempty"`
	Imageurl string             `json:"imageurl,omitempty" bson:"imageurl,omitempty"`
	Time     string             `json:"time,omitempty" bson:"time,omitempty"`
}

type postHandlers struct {
	sync.Mutex
	storep map[string]Post
}

func (h *postHandlers) posts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.getp(w, r)
		return
	case "POST":
		h.postp(w, r)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	}
}
func (h *postHandlers) getp(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	var posts []Post
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		posts = append(posts, post)
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(posts)
}
func (h *postHandlers) getPost(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, _ := primitive.ObjectIDFromHex(parts[2])
	var post Post
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := collection.FindOne(ctx, User{Id: id}).Decode(&post)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(post)
}
func (h *postHandlers) getPostuser(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	parts := strings.Split(r.URL.String(), "/")
	if len(parts) != 4 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id := parts[3]
	var posts []Post
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var post Post
		cursor.Decode(&post)
		if post.User == id {
			posts = append(posts, post)
		}
	}
	if err := cursor.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(posts)

}
func (h *postHandlers) postp(w http.ResponseWriter, r *http.Request) {
	r.Header.Add("content", "application/json")
	var post Post
	json.NewDecoder(r.Body).Decode(&post)
	collection := client.Database("test").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	post.Time = fmt.Sprintf(" %s", time.Unix(time.Now().Unix(), 0))
	result, _ := collection.InsertOne(ctx, post)
	json.NewEncoder(w).Encode(result)

}

func newpostHandlers() *postHandlers {
	return &postHandlers{
		storep: map[string]Post{},
	}
}
func main() {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))

	defer client.Disconnect(ctx)
	userHandlers := newuserHandlers()
	postHandlers := newpostHandlers()

	http.HandleFunc("/users", userHandlers.users)
	http.HandleFunc("/users/", userHandlers.getUser)
	http.HandleFunc("/posts", postHandlers.posts)
	http.HandleFunc("/posts/", postHandlers.getPost)
	http.HandleFunc("/posts/users/", postHandlers.getPostuser)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
