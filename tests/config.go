package tests

import (
	"reddit-activity-go-api/app"
	"reddit-activity-go-api/app/services"
)

type TestUtils struct {
	RedditService *services.RedditService
	RedditApp     *app.RedditApp
}

var Fixture TestUtils

func init() {
	Fixture.RedditService = services.NewRedditService()

	subReddits := []string{
		"AskReddit",
		"funny",
	}
	redditApp := app.InitApp(subReddits)
	Fixture.RedditApp = &redditApp
}
