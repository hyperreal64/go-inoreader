package userinfo

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const userInfoURL = "https://www.inoreader.com/reader/api/0/user-info"

// UserInfo JSON response
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

// Gets the user info. Sends a GET request and returns JSON response
// as a UserInfo struct.
func GetUserInfo(rc *resty.Client) (userinfo *UserInfo, err error) {
	resp, err := rc.R().Get(userInfoURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get user info")
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &userinfo); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", &userinfo)
	}

	return userinfo, nil
}
