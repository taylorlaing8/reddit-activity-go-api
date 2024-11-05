package app

import (
	// "reddit-activity-go-api/app/handlers"

	"cmp"
	"fmt"
	"net/http"
	"reddit-activity-go-api/app/models"
	"reddit-activity-go-api/app/services"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RedditApp struct {
	Router        *gin.Engine
	SubReddits    []string
	RedditService *services.RedditService
}

// Storing aggregation of posts in memory here
type PostsMeta struct {
	Posts       map[string]models.Post
	LastUpdated time.Time
	BeforeKey   *string
}

var PostsMap map[string]PostsMeta = make(map[string]PostsMeta)

func InitApp(subreddits []string) RedditApp {
	app := RedditApp{
		Router:        gin.Default(),
		SubReddits:    subreddits,
		RedditService: services.NewRedditService(),
	}

	app.setRoutes()

	// Set watch feeds to a go routine to run concurrently with api
	app.initFeeds()
	go app.watchFeeds()

	return app
}

func (app *RedditApp) setRoutes() {
	v1 := app.Router.Group("/v1")
	{
		v1Posts := v1.Group("/posts")
		{
			v1Posts.GET("/popular", func(c *gin.Context) {
				limit := 5

				limitStr := c.Query("limit")
				if len(limitStr) > 0 {
					limit, _ = strconv.Atoi(limitStr)
				}

				highestRatedPosts := app.GetHighestRatedPosts(limit)
				c.JSON(http.StatusOK, highestRatedPosts)
			})

			v1Posts.GET("", func(c *gin.Context) {
				allPosts := app.GetPosts()
				c.JSON(http.StatusOK, allPosts)
			})
		}

		v1Users := v1.Group("/users")
		{
			v1Users.GET("/engaged", func(c *gin.Context) {
				limit := 5

				limitStr := c.Query("limit")
				if len(limitStr) > 0 {
					limit, _ = strconv.Atoi(limitStr)
				}

				mostActiveUsers := app.GetMostActiveUsers(limit)
				c.JSON(http.StatusOK, mostActiveUsers)
			})
		}
	}
}

func (app *RedditApp) initFeeds() {
	for _, subreddit := range app.SubReddits {
		limit := 1
		subredditPostsRes, err := app.RedditService.GetNewPosts(subreddit, &limit, nil)
		if err != nil {
			fmt.Println(fmt.Errorf("error fetching subreddit posts: %v", err))
			continue
		}

		pMeta := PostsMeta{
			Posts:       make(map[string]models.Post),
			LastUpdated: time.Now(),
			BeforeKey:   nil,
		}

		if len(subredditPostsRes.Data.Children) > 0 {
			post := subredditPostsRes.Data.Children[0]

			// pMeta.Posts[post.Data.Name] = post
			pMeta.BeforeKey = &post.Data.Name
		}

		PostsMap[subreddit] = pMeta
	}
}

func (app *RedditApp) watchFeeds() {
	for {
		for _, subreddit := range app.SubReddits {
			postMeta := PostsMap[subreddit]

			limit := 100
			subredditPostsRes, err := app.RedditService.GetNewPosts(subreddit, &limit, postMeta.BeforeKey)
			if err != nil {
				fmt.Println(fmt.Errorf("error fetching subreddit posts: %v", err))
				continue
			}

			/*
			*  Note:
			*    - By not setting BeforeKey each time, system gains ability to track changes in UpVotes over time
			*    - However, this limits this simple program to only keep the last 100 posts updated (limiting calls to a single call with limit of 100 - Reddit's max)
			 */
			for _, postContent := range subredditPostsRes.Data.Children {
				postMeta.Posts[postContent.Data.Name] = postContent
			}

			postMeta.LastUpdated = time.Now()

			PostsMap[subreddit] = postMeta
		}

		time.Sleep(app.RedditService.HttpClient.BackoffSleepBase)
	}
}

func (app *RedditApp) GetPosts() []models.Post {
	allPosts := make([]models.Post, 0)

	for _, subredditPost := range PostsMap {
		for _, post := range subredditPost.Posts {
			allPosts = append(allPosts, post)
		}
	}

	slices.SortFunc(allPosts, func(a, b models.Post) int {
		return cmp.Compare(b.Data.Created, a.Data.Created)
	})

	return allPosts
}

func (app *RedditApp) GetHighestRatedPosts(limit int) []models.Post {
	// Another option: Store a heap that gets reprioritized asynchronously as posts are fetched. Then retrieving this list is O(limit) time (very efficient)
	orderedPosts := app.GetPosts()
	slices.SortFunc(orderedPosts, func(a, b models.Post) int {
		return cmp.Compare(b.Data.Ups, a.Data.Ups)
	})

	limit = min(limit, len(orderedPosts))
	return orderedPosts[:limit]
}

type AuthorPostCount struct {
	AuthorId  string `json:"authorId"`
	PostCount int    `json:"postCount"`
}

func (app *RedditApp) GetMostActiveUsers(limit int) []AuthorPostCount {
	// Another option: Store a heap that gets reprioritized asynchronously as posts are fetched. Then retrieving this list is O(limit) time (very efficient)
	allPosts := app.GetPosts()

	authorsPostCount := make(map[string]int)
	for _, post := range allPosts {
		if authorCount, ok := authorsPostCount[post.Data.Author]; ok {
			authorsPostCount[post.Data.Author] = authorCount + 1
		} else {
			authorsPostCount[post.Data.Author] = 1
		}
	}

	authorPostCountsList := make([]AuthorPostCount, 0)
	for author := range authorsPostCount {
		authorPostCountsList = append(authorPostCountsList, AuthorPostCount{
			AuthorId:  author,
			PostCount: authorsPostCount[author],
		})
	}

	slices.SortFunc(authorPostCountsList, func(a, b AuthorPostCount) int {
		return cmp.Compare(b.PostCount, a.PostCount)
	})

	limit = min(limit, len(authorPostCountsList))
	return authorPostCountsList[:limit]
}
