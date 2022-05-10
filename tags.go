package main

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
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

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), tfl); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", tfl)
	}

	return tfl, nil
}

// RenameTag --- Rename tag specified in query parameters
func RenameTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/rename-tag")
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
		Post(baseURL + "/edit-tag")
	if err != nil {
		return err
	}

	return nil
}

// TODO: Notify user when item cannot be marked unread due to `timestampUsec` being older than `firstitemsec` of its feed
// ExecEditTagRead -- Execute EditTag to mark the item as read (a) or unread (r)
func ExecEditTagRead(itemID string, markRead bool) error {

	var params = map[string]string{}
	if markRead {
		params = map[string]string{
			"a": "user/-/state/com.google/read",
			"i": itemID,
		}
	} else {
		params = map[string]string{
			"r": "user/-/state/com.google/read",
			"i": itemID,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as read", itemID)
	}

	return nil
}

// ExecEditTagStar --- Execute EditTag to star (a) or unstar (r) an item
func ExecEditTagStar(itemID string, starred bool) error {
	var params = map[string]string{}
	if starred {
		params = map[string]string{
			"a": "user/-/state/com.google/starred",
			"i": itemID,
		}
	} else {
		params = map[string]string{
			"r": "user/-/state/com.google/starred",
			"i": itemID,
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark item %s as %s", itemID, starred)
	}

	return nil
}

// func EditTagSaved(url string, saved bool) error {

// 	var params = map[string]string{}
// 	if saved {
// 		params = map[string]string{
// 			"a": "user/-/state/com.google/saved-web-pages",
// 			"i": url,
// 		}
// 	} else {
// 		params = map[string]string{
// 			"r": "user/-/state/com.google/saved-web-pages",
// 			"i": url,
// 		}
// 	}

// 	ctx, cancel := context.WithCancel(context.Background())
// 	rClient := oauth2RestyClient(ctx)
// 	defer cancel()

// 	if err := EditTag(rClient, params); err != nil {
// 		return errors.Wrapf(err, "Could not mark item %s as %s", url, saved)
// 	}

// 	return
// }

// ExecRenameTag --- Execute EditTag to rename a tag from src to dest
func ExecRenameTag(src string, dest string) error {

	params := map[string]string{
		"s":    src,
		"dest": dest,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := RenameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not rename tag %s to %s", src, dest)
	}

	return nil
}

// ExecDelTag --- Execute EditTag to delete a tag
func ExecDelTag(tagName string) error {

	params := map[string]string{"s": tagName}

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	if err := RenameTag(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not delete tag %s", tagName)
	}

	return nil
}
