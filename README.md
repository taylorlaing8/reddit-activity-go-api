# Reddit Activity Simple API

Reddit, much like other social media platforms, provides a way for users to commmunicate their interests. This simple application allows following (watching) a list of subreddits that may interest the user. From the time the application starts, new posts will be tracked, and the most recent 100 posts will be kept up to date in terms of content (such as Upvote count).

Data is fetched on a separate thread (go routine) to prevent conjestion when interacting with this application. All fetches to the Reddit API use the response headers to calculate how long to wait between calls. This allows for fetching the most recent data as frequently as possible, without causing rate limit errors. However, should an error occur, retries with an incremental backoff will be hit to attempt to get the data still.

## Running the Application
Assuming your machine has Go installed, you can simply run `go run .` within the terminal while at the root of this project. This will start the HTTP server on port 8080, as specified within the `main.go` file. Once this has started, requests to the Reddit API will automatically begin in the background to start building the list of posts that are created after application start time. From there, API requests can be made (using the endpoints below) to retrieve data that has been located from this background activity. This method ensures that proper rate limiting is maintained with the Reddit API, while not passing on any rate limiting to the endpoints below.

## Testing the Application
To run all tests, which go through unit testing the data models, and ensuring correct responses from the Reddit Service and App, simply run `go test ./tests` to kick off all test runs. Depending on your IDE, these tests could also be run individually as well (such as in VS Code).

## API Endpoints
### Reddit Posts
`GET http://localhost:8080/posts`
Returns the full list of posts that were created within the time of app start, up to the time of this GET request. Posts are ordered by creation time, and contain a mix of all posts from any subreddit that is being monitored.

`GET http://localhost:8080/v1/posts/popular?limit=<LIMIT>`
Returns the most popular posts as determined by their respective UpVote count. Optional query parameter of `limit` can be given to return a subset of these ordered posts.

### Users
`GET http://localhost:8080/v1/users/engaged?limit=<LIMIT>`
Returns the most engaged users as determined by the count of their posts across all monitored subreddits. Optional query parameter of `limit` can be given to return a subset of these ordered posts. Note that given the likely short-lived nature of this app, it is unlikely that there will be more than 1 post per user.