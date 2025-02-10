package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/TimiBolu/lema-ai-users-service/config"
	"github.com/TimiBolu/lema-ai-users-service/database"
	"github.com/TimiBolu/lema-ai-users-service/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func GetPostsByUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		http.Error(w, "userId parameter is required", http.StatusBadRequest)
		return
	}

	var posts []models.Post
	result := database.DB.Where("user_id = ?", userID).Find(&posts)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	var postData struct {
		Title  string `json:"title"`
		Body   string `json:"body"`
		UserID string `json:"userId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&postData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if postData.Title == "" || postData.Body == "" || postData.UserID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := database.DB.First(&user, "id = ?", postData.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	newPost := models.Post{
		ID:     uuid.NewString(),
		Title:  postData.Title,
		Body:   postData.Body,
		UserID: postData.UserID,
	}

	if err := database.DB.Create(&newPost).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// SQLite is quite fast so animations don't show locally
	if config.EnvConfig.APP_ENV == config.ServerEnvironmentDevelopment {
		// Simulate network delays
		// Slow the API down by 1 second
		time.Sleep(1 * time.Second)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newPost)
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var post models.Post
	result := database.DB.First(&post, "id = ?", id)
	if result.Error != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	database.DB.Delete(&post)
	w.WriteHeader(http.StatusNoContent)
}
