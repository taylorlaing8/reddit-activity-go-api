package models

import (
	"encoding/json"
	"fmt"
)

type PostType int

const (
	Comment PostType = iota
	Account
	Link
	Message
	Subreddit
	Award
)

func (pType PostType) String() string {
	return [...]string{
		"t1",
		"t2",
		"t3",
		"t4",
		"t5",
		"t6",
	}[pType]
}

func (pType PostType) MarshalJSON() ([]byte, error) {
	return json.Marshal(pType.String())
}

func (pType *PostType) UnmarshalJSON(data []byte) error {
	var accountTypeString string = ""

	err := json.Unmarshal(data, &accountTypeString)
	if err != nil {
		return err
	}

	pType, err = GetPostType(&accountTypeString)

	return err
}

func GetPostType(pType *string) (*PostType, error) {
	var matchedPostType PostType

	if pType == nil || len(*pType) == 0 {
		return nil, fmt.Errorf("cannot get post type for nil or empty string")
	} else {
		switch *pType {
		case Comment.String():
			matchedPostType = Comment
		case Account.String():
			matchedPostType = Account
		case Link.String():
			matchedPostType = Link
		case Message.String():
			matchedPostType = Message
		case Subreddit.String():
			matchedPostType = Subreddit
		case Award.String():
			matchedPostType = Award
		default:
			return nil, fmt.Errorf("no matching account type found for: %v", pType)
		}
	}

	return &matchedPostType, nil
}

type Post struct {
	Kind PostType    `json:"kind"`
	Data PostContent `json:"data"`
}

type PostContent struct {
	PostID            string  `json:"id"`
	Name              string  `json:"name"` // Used for 'after' value in pagination
	Title             string  `json:"title"`
	SubReddit         string  `json:"subreddit"`
	SubRedditPrefixed string  `json:"subreddit_name_prefixed"`
	SubRedditId       string  `json:"subreddit_id"`
	Downs             int     `json:"downs"`
	Ups               int     `json:"ups"`
	UpvoteRatio       float32 `json:"upvote_ratio"`
	RobotIndexable    bool    `json:"is_robot_indexable"`
	Author            string  `json:"author"`
	URL               string  `json:"url"`
	Created           float32 `json:"created"`
}
