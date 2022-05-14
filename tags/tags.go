package tags

import (
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

const (
	tagListURL   = "https://www.inoreader.com/reader/api/0/tag/list?types=1&counts=1"
	renameTagURL = "https://www.inoreader.com/reader/api/0/rename-tag"
	editTagURL   = "https://www.inoreader.com/reader/api/0/edit-tag"
)

// TagFolderList response
type TagFolderList struct {
	Tags []struct {
		ID          string `json:"id"`
		Sortid      string `json:"sortid"`
		UnreadCount int64  `json:"unread_count"`
		Type        string `json:"type"`
	} `json:"tags"`
}

// GetTagList --- Get list of tags
func GetTagList(rc *resty.Client) (tfl *TagFolderList, err error) {

	resp, err := rc.R().Get(tagListURL)
	if err != nil {
		return nil, err
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &tfl); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", tfl)
	}

	return tfl, nil
}

// RenameTag --- Rename tag specified in query parameters
func RenameTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(renameTagURL)
	if err != nil {
		return err
	}

	return nil
}

// func DeleteTag(rc *resty.Client, tagName string) error {

// 	_, err := rc.R().
// 		SetQueryParams(map[string]string{
// 			"s": tagName,
// 		}).
// 		Post(baseURL + "/disable-tag")
// 	if err != nil {
// 		return err
// 	}

// 	return
// }

// EditTag -- Edit tag specified in query parameters
func EditTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(editTagURL)
	if err != nil {
		return err
	}

	return nil
}
