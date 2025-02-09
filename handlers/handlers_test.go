package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TimiBolu/lema-ai-users-service/database"
	"github.com/TimiBolu/lema-ai-users-service/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestSuite struct to manage setup and teardown
type HandlerTestSuite struct {
	suite.Suite
	DB *gorm.DB
}

// SetupSuite initializes the test database
func (suite *HandlerTestSuite) SetupSuite() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(suite.T(), err)
	suite.DB = db
	database.DB = db

	// Auto-migrate tables
	suite.DB.AutoMigrate(&models.User{}, &models.Address{}, &models.Post{})

	// Insert mock data
	user := models.User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john@example.com"}
	suite.DB.Create(&user)

	address := models.Address{ID: "1", UserID: "1", Street: "123 Main St", City: "Lagos", State: "LA", ZipCode: "100001"}
	suite.DB.Create(&address)

	post := models.Post{ID: "1", UserID: "1", Title: "First Post", Body: "This is a test post"}
	suite.DB.Create(&post)
}

// Test GetUsers
func (suite *HandlerTestSuite) TestGetUsers() {
	req, _ := http.NewRequest("GET", "/users?pageNumber=1&pageSize=10", nil)
	rr := httptest.NewRecorder()

	GetUsers(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var response map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Contains(suite.T(), response, "users")
	assert.Contains(suite.T(), response, "pagination")
}

// Test GetUserByID
func (suite *HandlerTestSuite) TestGetUserByID() {
	req, _ := http.NewRequest("GET", "/users/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", GetUserByID)
	router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var user models.User
	json.Unmarshal(rr.Body.Bytes(), &user)
	assert.Equal(suite.T(), "1", user.ID)
}

// Test GetUsersCount
func (suite *HandlerTestSuite) TestGetUsersCount() {
	req, _ := http.NewRequest("GET", "/users/count", nil)
	rr := httptest.NewRecorder()

	GetUsersCount(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var response map[string]int64
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(suite.T(), int64(1), response["count"])
}

// Test GetPostsByUser
func (suite *HandlerTestSuite) TestGetPostsByUser() {
	req, _ := http.NewRequest("GET", "/posts?userId=1", nil)
	rr := httptest.NewRecorder()

	GetPostsByUser(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var posts []models.Post
	json.Unmarshal(rr.Body.Bytes(), &posts)
	assert.NotEmpty(suite.T(), posts)
	assert.Equal(suite.T(), "1", posts[0].UserID)
}

// Test CreatePost
func (suite *HandlerTestSuite) TestCreatePost() {
	postData := `{"title": "New Post", "body": "This is a new post", "userId": "1"}`
	req, _ := http.NewRequest("POST", "/posts", bytes.NewBuffer([]byte(postData)))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	CreatePost(rr, req)

	assert.Equal(suite.T(), http.StatusCreated, rr.Code)

	var post models.Post
	json.Unmarshal(rr.Body.Bytes(), &post)
	assert.Equal(suite.T(), "New Post", post.Title)
}

// Test DeletePost
func (suite *HandlerTestSuite) TestDeletePost() {
	req, _ := http.NewRequest("DELETE", "/posts/1", nil)
	rr := httptest.NewRecorder()

	router := mux.NewRouter()
	router.HandleFunc("/posts/{id}", DeletePost)
	router.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusNoContent, rr.Code)

	var post models.Post
	result := suite.DB.First(&post, "id = ?", "1")
	assert.Error(suite.T(), result.Error) // Should not find post
}

// Run tests
func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
