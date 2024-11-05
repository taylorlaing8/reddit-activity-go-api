package services

import (
	"fmt"
	"strconv"
	"time"

	"reddit-activity-go-api/app/models"
)

type RedditService struct {
	baseUrl            string
	HttpClient         *HttpClient
	rateLimitUsed      int
	rateLimitRemaining float64
	rateLimitReset     int
}

func NewRedditService() *RedditService {
	client := NewClient(
		15*time.Second,
		[]time.Duration{
			0 * time.Millisecond,
			500 * time.Millisecond,
			1500 * time.Millisecond,
			3000 * time.Millisecond,
		},
	)

	return &RedditService{
		baseUrl:            "https://www.reddit.com",
		HttpClient:         client,
		rateLimitUsed:      0,
		rateLimitRemaining: 1,
		rateLimitReset:     0,
	}
}

func (r *RedditService) GetRateLimits() (int, float64, int) {
	return r.rateLimitUsed, r.rateLimitRemaining, r.rateLimitReset
}

func (r *RedditService) CalculateRateLimitSleep() {
	_, rateLimitRemaining, rateLimitReset := r.GetRateLimits()
	timeoutPeriod := time.Duration(1000 * float64(rateLimitReset) / rateLimitRemaining)
	r.HttpClient.BackoffSleepBase = (timeoutPeriod * time.Millisecond)
}

func (r *RedditService) GetNewPosts(subreddit string, limit *int, before *string) (*models.Listing, error) {
	pageLimit := 25
	if limit != nil {
		pageLimit = *limit
	}

	apiUrl := fmt.Sprintf("%s/r/%s/new.json?limit=%d", r.baseUrl, subreddit, pageLimit)
	if before != nil {
		apiUrl += fmt.Sprintf("&before=%s", *before)
	}

	responseBody, responseHeader, err := SendGet[models.Listing](*r.HttpClient, apiUrl)
	if err != nil {
		return nil, err
	}

	r.CalculateRateLimitSleep()

	if responseHeader != nil {
		header := *responseHeader
		if rateLimitUsed, ok := header["X-Ratelimit-Used"]; ok && len(rateLimitUsed) > 0 {
			r.rateLimitUsed, _ = strconv.Atoi(rateLimitUsed[0])
		}
		if rateLimitRemaining, ok := header["X-Ratelimit-Remaining"]; ok && len(rateLimitRemaining) > 0 {
			r.rateLimitRemaining, _ = strconv.ParseFloat(rateLimitRemaining[0], 64)
		}
		if rateLimitReset, ok := header["X-Ratelimit-Reset"]; ok && len(rateLimitReset) > 0 {
			r.rateLimitReset, _ = strconv.Atoi(rateLimitReset[0])
		}
	}

	return responseBody, nil
}
