package tests

import (
	"encoding/json"
	"flag"
	"net/http"
	"net/http/httptest"
	"reddit-activity-go-api/app/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Reddit_Service_Fetch_Posts_Should_Succeed(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/posts", nil)
	Fixture.RedditApp.Router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	responseBody := w.Body.String()

	var posts []models.Post

	err := json.Unmarshal([]byte(responseBody), &posts)
	assert.Nil(t, err)
}

func Test_Reddit_Service_Fetch_Posts_Should_Succeed_And_Return_Non_Zero_Response_When_Waiting(t *testing.T) {
	flag.Set("test.timeout", "120s")

	var posts []models.Post

	// Loop until a post is returned
	retryCount := 0
	for retryCount < 90 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/posts", nil)
		Fixture.RedditApp.Router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		responseBody := w.Body.String()

		err := json.Unmarshal([]byte(responseBody), &posts)
		assert.Nil(t, err)

		if len(posts) > 0 {
			assert.NotNil(t, posts[0].Data)
			break
		} else {
			retryCount += 1
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func Test_Reddit_Service_Fetch_Popular_Posts_Should_Succeed_And_Return_Non_Zero_Response_When_Waiting(t *testing.T) {
	flag.Set("test.timeout", "120s")

	var posts []models.Post

	// Loop until a post is returned
	retryCount := 0
	for retryCount < 90 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/posts/popular", nil)
		Fixture.RedditApp.Router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		responseBody := w.Body.String()

		err := json.Unmarshal([]byte(responseBody), &posts)
		assert.Nil(t, err)

		if len(posts) > 0 {
			assert.NotNil(t, posts[0].Data)
			break
		} else {
			retryCount += 1
			time.Sleep(2000 * time.Millisecond)
		}
	}
}

func Test_Reddit_Service_Fetch_Engaged_Users_Should_Succeed_And_Return_Non_Zero_Response_When_Waiting(t *testing.T) {
	flag.Set("test.timeout", "120s")

	var posts []models.Post

	// Loop until a post is returned
	retryCount := 0
	for retryCount < 90 {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/v1/users/engaged", nil)
		Fixture.RedditApp.Router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		responseBody := w.Body.String()

		err := json.Unmarshal([]byte(responseBody), &posts)
		assert.Nil(t, err)

		if len(posts) > 0 {
			assert.NotNil(t, posts[0].Data)
			break
		} else {
			retryCount += 1
			time.Sleep(2000 * time.Millisecond)
		}
	}
}
