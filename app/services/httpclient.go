package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	httpClient       *http.Client
	BackoffSleepBase time.Duration
	backoffSchedule  []time.Duration
}

func NewClient(timeout time.Duration, backoffSchedule []time.Duration) *HttpClient {
	return &HttpClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		backoffSchedule:  backoffSchedule,
		BackoffSleepBase: 1000 * time.Millisecond,
	}
}

func SendGet[TResponse interface{}](client HttpClient, apiUrl string) (res *TResponse, header *http.Header, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	httpResponse, responseHeader, err := sendRequestWithoutBody[TResponse](client, "GET", apiUrl)
	if err != nil {
		return nil, nil, err
	}

	return httpResponse, responseHeader, nil
}

func sendRequestWithoutBody[TResponse interface{}](client HttpClient, httpMethod string, url string) (*TResponse, *http.Header, error) {
	httpRequest := createRequest(httpMethod, url, nil)

	httpResponse, responseHeader, err := handleBackoffRequest[TResponse](client, httpRequest)
	if err != nil {
		return nil, nil, err
	}

	return httpResponse, responseHeader, nil
}

func createRequest(httpMethod string, url string, requestBody *[]byte) *http.Request {
	var httpRequest *http.Request

	// Set up to handle with and without body, even though this sample app only ever uses without (i.e. GET requests)
	if requestBody == nil {
		httpRequest, _ = http.NewRequest(httpMethod, url, nil)
	} else {
		httpRequest, _ = http.NewRequest(httpMethod, url, bytes.NewBuffer(*requestBody))
	}

	httpRequest.Header.Set("Accept", "application/json")

	return httpRequest
}

func handleBackoffRequest[TResponse interface{}](client HttpClient, httpRequest *http.Request) (*TResponse, *http.Header, error) {
	var httpResponse *TResponse
	var responseHeader *http.Header
	var err error

	for _, backoff := range client.backoffSchedule {
		httpResponse, responseHeader, err = handleRequest[TResponse](client, httpRequest)
		if err == nil {
			break
		}

		time.Sleep(client.BackoffSleepBase + backoff)
	}

	if err != nil {
		return nil, nil, err
	}

	return httpResponse, responseHeader, nil
}

func handleRequest[TResponse interface{}](client HttpClient, httpRequest *http.Request) (*TResponse, *http.Header, error) {
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("error fetching new posts: %v", err)
	}
	defer httpResponse.Body.Close()

	content, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("error parsing response body: %v", err)
	}

	switch httpResponse.StatusCode {
	case http.StatusOK:
		var result TResponse
		json.Unmarshal(content, &result)

		return &result, &httpResponse.Header, nil
	default:
		return nil, nil, fmt.Errorf("status code: %v | response:\n%s", httpResponse.StatusCode, content)
	}
}
