package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Reddit_App_Fetch_Posts_Should_Succeed(t *testing.T) {
	subReddit := "AskReddit"
	limit := 5

	res, err := Fixture.RedditService.GetNewPosts(subReddit, &limit, nil)
	assert.Nil(t, err)
	assert.NotNil(t, res.Data)
	assert.Len(t, res.Data.Children, 5)
}

func Test_Reddit_App_Fetch_Posts_Should_Succeed_And_Return_Earlier_Results(t *testing.T) {
	subReddit := "AskReddit"
	limit := 5

	res, err := Fixture.RedditService.GetNewPosts(subReddit, &limit, nil)
	assert.Nil(t, err)
	assert.NotNil(t, res.Data)
	assert.Len(t, res.Data.Children, 5)

	tokenOne := res.Data.Children[0].Data.Name
	tokenTwo := res.Data.Children[1].Data.Name
	res2, err := Fixture.RedditService.GetNewPosts(subReddit, &limit, &tokenTwo)
	assert.Nil(t, err)
	assert.NotNil(t, res2.Data)
	assert.Len(t, res2.Data.Children, 1)
	assert.Equal(t, tokenOne, res2.Data.Children[0].Data.Name)
}
