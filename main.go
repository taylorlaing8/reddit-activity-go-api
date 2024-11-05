package main

import (
	"reddit-activity-go-api/app"
)

func main() {
	subReddits := []string{
		"AskReddit",
		"funny",
	}

	redditApp := app.InitApp(subReddits)

	redditApp.Router.Run(":8080")
}
