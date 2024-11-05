package tests

import (
	"reddit-activity-go-api/app/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PostType_Match_Succeeds_For_Valid_Strings(t *testing.T) {
	postTypes := []models.PostType{
		models.Comment,
		models.Account,
		models.Link,
		models.Message,
		models.Subreddit,
		models.Award,
	}

	for _, postType := range postTypes {
		postTypeString := postType.String()

		locatedPostType, err := models.GetPostType(&postTypeString)
		assert.Nil(t, err)
		assert.Equal(t, postType, *locatedPostType)
	}
}

func Test_PostType_Match_Should_Have_Errors_For_Invalid_Strings(t *testing.T) {
	postTypes := []string{"NoPostType", "Invalid", "Fake"}

	for _, postType := range postTypes {
		locatedPostType, err := models.GetPostType(&postType)
		assert.Contains(t, err.Error(), "no matching account type found for")
		assert.Nil(t, locatedPostType)
	}
}
