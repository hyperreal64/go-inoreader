package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
)

const (
	baseURL           = "https://www.inoreader.com/reader/api/0"
	userInfoURL       = baseURL + "/user-info"
	addSubURL         = baseURL + "/subscription/quickadd"
	editSubURL        = baseURL + "/subscription/edit"
	unreadCountersURL = baseURL + "/unread-count"
	subListURL        = baseURL + "/subscription/list"
	tagListURL        = baseURL + "/tag/list"
	streamContentURL  = baseURL + "/stream/contents/"
	itemIDsURL        = baseURL + "/stream/items/ids"
	streamPrefsURL    = baseURL + "/preference/stream/list"
	streamPrefsSetURL = baseURL + "/preference/stream/set"
)

func httpDo(client *http.Client, method string, url string) ([]byte, error) {

	var (
		res *http.Response
		err error
	)

	if method == "GET" {
		res, err = client.Get(url)
	}

	if method == "POST" {
		res, err = client.Post(url, "application/json", nil)
	}

	if err != nil {
		return nil, errors.Wrap(err, "Could not complete HTTP request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read HTTP response body")
	}

	return body, nil
}

// SetInoreader --- Makes changes to the Inoreader user's account.
func SetInoreader(client *http.Client, url string, params interface{}) error {

	v, err := query.Values(params)
	if err != nil {
		return errors.Wrapf(err, "Could not construct URL with query parameters: %v", params)
	}

	encodedURL := fmt.Sprintf("%s?%s", url, v.Encode())

	_, err = httpDo(client, "POST", encodedURL)
	if err != nil {
		return err
	}

	return nil
}

// GetUserInfo ---
func GetUserInfo(client *http.Client, userInfo *UserInfo) error {

	body, err := httpDo(client, "GET", userInfoURL)
	if err != nil {
		return errors.Wrap(err, "Could not get user info")
	}

	if err := json.Unmarshal(body, userInfo); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", userInfo)
	}

	return nil
}

// QuickAddSubscription ---
func QuickAddSubscription(client *http.Client, params *QuickAddParams) error {

	v, err := query.Values(params)
	if err != nil {
		return errors.Wrapf(err, "Could not construct URL with query parameters: %v", params)
	}

	encodedURL := fmt.Sprintf("%s?%s", addSubURL, v.Encode())

	body, err := httpDo(client, "POST", encodedURL)
	if err != nil {
		return errors.Wrap(err, "Could not add subscription")
	}

	quickAdd := &QuickAdd{}
	if err := json.Unmarshal(body, quickAdd); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", quickAdd)
	}

	if quickAdd.NumResults != 1 {
		return errors.New("Feed not added")
	}

	return nil
}

// EditSubscription ---
func EditSubscription(client *http.Client, params *EditSubParams) error {

	if err := SetInoreader(client, editSubURL, params); err != nil {
		return errors.Wrap(err, "Could not edit subscription")
	}

	return nil
}

// GetUnreadCounters ---
func GetUnreadCounters(client *http.Client, unreadCounters *UnreadCounters) error {

	body, err := httpDo(client, "GET", unreadCountersURL)
	if err != nil {
		return errors.Wrap(err, "Could not get unread counters")
	}

	if err = json.Unmarshal(body, &unreadCounters); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object %v", unreadCounters)
	}

	return nil
}

// GetSubscriptionList ---
func GetSubscriptionList(client *http.Client, subList *SubscriptionList) error {

	body, err := httpDo(client, "GET", subListURL)
	if err != nil {
		return errors.Wrap(err, "Could not get subscription list")
	}

	if err := json.Unmarshal(body, subList); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", subList)
	}

	return nil
}

// GetTagList ---
func GetTagList(client *http.Client, tagList *TagFolderList) error {

	body, err := httpDo(client, "GET", tagListURL)
	if err != nil {
		return errors.Wrap(err, "Could not get tag/folder list")
	}

	if err := json.Unmarshal(body, tagList); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", tagList)
	}

	return nil
}

// GetStreamContents ---
func GetStreamContents(client *http.Client, streamContents *StreamContents, params *ContentsParams) error {

	v, err := query.Values(params)
	if err != nil {
		return errors.Wrapf(err, "Could not construct URL with query parameters: %v", params)
	}

	encodedURL := fmt.Sprintf("%s?%s", streamContentURL, v.Encode())

	body, err := httpDo(client, "GET", encodedURL)
	if err != nil {
		return errors.Wrap(err, "Could not get stream contents")
	}

	if err := json.Unmarshal(body, streamContents); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", streamContents)
	}

	return nil
}

// GetItemIDs ---
func GetItemIDs(client *http.Client, itemIDs *ItemIDs, params *ContentsParams) error {

	v, err := query.Values(params)
	if err != nil {
		return errors.Wrapf(err, "Could not construct URL with query parameters: %v", params)
	}

	encodedURL := fmt.Sprintf("%s?%s", itemIDsURL, v.Encode())

	body, err := httpDo(client, "GET", encodedURL)
	if err != nil {
		return errors.Wrap(err, "Could not get item IDs")
	}

	if err := json.Unmarshal(body, itemIDs); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", itemIDs)
	}

	return nil
}

// GetStreamPrefsList ---
func GetStreamPrefsList(client *http.Client, streamPrefsList *StreamPreferenceList) error {

	body, err := httpDo(client, "GET", streamPrefsURL)
	if err != nil {
		return errors.Wrap(err, "Could not get stream preferences list")
	}

	if err := json.Unmarshal(body, streamPrefsList); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", streamPrefsList)
	}

	return nil
}

// SetStreamPrefs ---
func SetStreamPrefs(client *http.Client, params *SetStreamPrefsParams) error {

	if err := SetInoreader(client, streamPrefsSetURL, params); err != nil {
		return errors.Wrap(err, "Could not set stream preferences")
	}

	return nil
}

// RenameTag ---
func RenameTag(client *http.Client, params *RenameTagParams) error {

	url := baseURL + "/rename-tag"
	if err := SetInoreader(client, url, params); err != nil {
		return errors.Wrapf(err, "Could not rename %s tag to %s", params.Source, params.Dest)
	}

	return nil
}

// DeleteTag ---
func DeleteTag(client *http.Client, tagName string) error {

	params := &DeleteTagParams{
		StreamID: tagName,
	}

	url := baseURL + "/disable-tag"
	if err := SetInoreader(client, url, params); err != nil {
		return errors.Wrapf(err, "Could not delete tag %s", tagName)
	}

	return nil
}

// EditTag ---
func EditTag(client *http.Client, params *EditTagParams) error {

	url := baseURL + "/edit-tag"
	if err := SetInoreader(client, url, params); err != nil {
		return errors.Wrap(err, "Could not edit tag")
	}

	return nil
}

// MarkAllAsRead ---
func MarkAllAsRead(client *http.Client, params *MarkAllAsReadParams) error {

	url := baseURL + "/mark-all-as-read"
	if err := SetInoreader(client, url, params); err != nil {
		return errors.Wrapf(err, "Could not mark all items in %s as read", params.StreamID)
	}

	return nil
}

// CreateActiveSearch ---
// Not implemented due to user account limitations

// DeleteActiveSearch ---
// Not implemented due to user account limitations
