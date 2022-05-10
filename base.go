package main

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

const (
	baseURL           = "https://www.inoreader.com/reader/api/0"
	userInfoURL       = baseURL + "/user-info"
	addSubURL         = baseURL + "/subscription/quickadd"
	editSubURL        = baseURL + "/subscription/edit"
	unreadCountersURL = baseURL + "/unread-count"
	subListURL        = baseURL + "/subscription/list"
	tagListURL        = baseURL + "/tag/list?types=1&counts=1"
	streamContentsURL = baseURL + "/stream/contents"
	itemIDsURL        = baseURL + "/stream/items/ids"
	streamPrefsURL    = baseURL + "/preference/stream/list"
	streamPrefsSetURL = baseURL + "/preference/stream/set"
)

// UserInfo response
// Output looks like:
// {
//     "userId": "1005869311",
//     "userName": "hyperreal",
//     "userProfileId": "1005869311",
//     "userEmail": "serio.jeffrey@gmail.com",
//     "isBloggerUser": false,
//     "signupTimeSec": 1379693893,
//     "isMultiLoginEnabled": false
// }
type UserInfo struct {
	UserID              string `json:"userId"`
	UserName            string `json:"userName"`
	UserProfileID       string `json:"userProfileId"`
	UserEmail           string `json:"userEmail"`
	IsBloggerUser       bool   `json:"isBloggerUser"`
	SignupTimeSec       int64  `json:"signupTimeSec"`
	IsMultiLoginEnabled bool   `json:"isMultiLoginEnabled"`
}

// GetUserInfo --- Gets the user info
func GetUserInfo() (userinfo *UserInfo, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	resp, err := rc.R().Get(userInfoURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get user info")
	}

	if err := json.Unmarshal(resp.Body(), &userinfo); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", &userinfo)
	}

	return nil, nil
}
