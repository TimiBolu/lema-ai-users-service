package handlers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/TimiBolu/lema-ai-users-service/config"
	"github.com/TimiBolu/lema-ai-users-service/database"
	"github.com/TimiBolu/lema-ai-users-service/models"
	"github.com/gorilla/mux"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	pageNumber, err1 := strconv.Atoi(r.URL.Query().Get("pageNumber"))
	pageSize, err2 := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if err1 != nil || pageNumber < 1 {
		pageNumber = 1
	}
	if err2 != nil || pageSize < 1 {
		pageSize = 10
	}

	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	var users []models.User
	result := database.DB.Preload("Address").
		Offset((pageNumber - 1) * pageSize).
		Limit(pageSize).
		Find(&users)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	totalPages := int(math.Ceil(float64(totalUsers) / float64(pageSize)))
	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"currentPage": pageNumber,
			"pageSize":    pageSize,
			"totalPages":  totalPages,
			"totalItems":  totalUsers,
			"hasNext":     pageNumber < totalPages,
			"hasPrev":     pageNumber > 1,
		},
	}

	// SQLite is quite fast so animations don't show locally
	if config.EnvConfig.APP_ENV == config.ServerEnvironmentDevelopment {
		// Simulate network delays
		// Slow the API down by 1 second
		time.Sleep(1 * time.Second)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user models.User
	result := database.DB.Preload("Address").First(&user, "id = ?", id)

	if result.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func GetUsersCount(w http.ResponseWriter, r *http.Request) {
	var count int64
	result := database.DB.Model(&models.User{}).Count(&count)

	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{"count": count})
}
