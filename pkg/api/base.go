package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

var authToken = ""

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

func httpRequest(apiURL string, method string, params map[string]string) ([]byte, error) {

	v := url.Values{}
	if len(params) != 0 {
		for key, value := range params {
			v.Add(key, value)
		}
	}

	url := fmt.Sprintf("%s?%s", apiURL, v.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "HTTP request failed")
	}

	req.Header.Add("Authorization", authToken)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to complete HTTP request")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read HTTP response body")
	}

	return body, nil
}

// GetJSONObject ---
func GetJSONObject(url string, req string, obj interface{}) error {

	body, err := httpRequest(url, req, nil)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &obj); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal JSON object %v", obj)
	}

	return nil
}

// GetJSONObjectParams ---
func GetJSONObjectParams(url string, req string, obj interface{}, params map[string]string) error {

	FilterParams(params)
	body, err := httpRequest(url, req, params)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &obj); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal JSON object %v", obj)
	}

	return nil
}

// FilterParams --- Filters unneeded url query params from the map provided.
// Maps are reference types like pointers and slices, so a map object passed
// to this function will be changed upon leaving this function's scope
func FilterParams(paramsMap map[string]string) {

	for key, value := range paramsMap {
		if len(value) == 0 {
			delete(paramsMap, key)
		}
	}
}

// SetInoreader --- Makes changes to the Inoreader user's account.
func SetInoreader(url string, paramsMap map[string]string) error {

	FilterParams(paramsMap)

	_, err := httpRequest(url, "POST", paramsMap)
	if err != nil {
		return err
	}

	return nil
}

// GetUserInfo ---
func GetUserInfo(userInfo *UserInfo) error {

	if err := GetJSONObject(userInfoURL, "GET", userInfo); err != nil {
		return errors.Wrap(err, "Failed to get user info")
	}

	return nil
}

// QuickAddSubscription ---
func QuickAddSubscription(feed string) error {

	params := map[string]string{
		"quickadd": feed,
	}

	quickAdd := &QuickAdd{}
	if err := GetJSONObjectParams(addSubURL, "POST", quickAdd, params); err != nil {
		return errors.Wrap(err, "Failed to add subscription")
	}

	if quickAdd.NumResults != 1 {
		return errors.New("Feed not added")
	}

	return nil
}

// EditSubscription ---
// when calling this function, the order of the url query params
// should be enforced so as to avoid mis-indexing. Also, len(params)
// should always == 5 to avoid index out of range error.
func EditSubscription(params []string) error {

	paramsMap := map[string]string{
		"ac": params[0],
		"s":  params[1],
		"t":  params[2],
		"a":  params[3],
		"r":  params[4],
	}

	if err := SetInoreader(editSubURL, paramsMap); err != nil {
		return errors.Wrap(err, "Failed to edit subscription")
	}

	return nil
}

// GetUnreadCounters ---
func GetUnreadCounters(unreadCounters *UnreadCounters) error {

	if err := GetJSONObject(unreadCountersURL, "GET", unreadCounters); err != nil {
		return errors.Wrap(err, "Failed to get unread counters")
	}

	return nil
}

// GetSubscriptionList ---
func GetSubscriptionList(subList *SubscriptionList) error {

	if err := GetJSONObject(subListURL, "GET", subList); err != nil {
		return errors.Wrap(err, "Failed to get subscription list")
	}

	return nil
}

// GetTagList ---
func GetTagList(tagList *TagFolderList) error {

	if err := GetJSONObject(tagListURL, "GET", tagList); err != nil {
		return errors.Wrap(err, "Failed to get tag/folder list")
	}

	return nil
}

// GetStreamContents ---
// Order of params should be enforced by calling function
func GetStreamContents(streamContents *StreamContents, params []string) error {

	paramsMap := map[string]string{
		"n":                         params[0],
		"r":                         params[1],
		"ot":                        params[2],
		"xt":                        params[3],
		"it":                        params[4],
		"c":                         params[5],
		"output":                    params[6],
		"includeAllDirectStreamIDs": params[7],
		"streamID":                  params[8],
	}

	if err := GetJSONObjectParams(streamContentURL, "GET", streamContents, paramsMap); err != nil {
		return errors.Wrap(err, "Failed to get stream contents")
	}

	return nil
}

// GetItemIDs ---
// Order of params should be enforced by calling function
func GetItemIDs(itemIDs *ItemIDs, params []string) error {

	paramsMap := map[string]string{
		"n":                         params[0],
		"r":                         params[1],
		"ot":                        params[2],
		"xt":                        params[3],
		"it":                        params[4],
		"c":                         params[5],
		"output":                    params[6],
		"IncludeAllDirectStreamIDs": params[7],
		"StreamID":                  params[8],
	}

	if err := SetInoreader(itemIDsURL, paramsMap); err != nil {
		return errors.Wrap(err, "Failed to get item ids")
	}

	return nil
}

// GetStreamPrefsList ---
func GetStreamPrefsList(streamPrefsList *StreamPreferenceList) error {

	if err := GetJSONObject(streamPrefsURL, "GET", streamPrefsList); err != nil {
		return errors.Wrap(err, "Failed to get stream preferences list")
	}

	return nil
}

// SetStreamPrefs ---
// Order of params should be enforced by calling function
func SetStreamPrefs(streamPrefParams []string) error {

	paramsMap := map[string]string{
		"s": streamPrefParams[0],
		"k": streamPrefParams[1],
		"v": streamPrefParams[2],
	}

	if err := SetInoreader(streamPrefsSetURL, paramsMap); err != nil {
		return errors.Wrap(err, "Failed to set stream preferences")
	}

	return nil
}

// RenameTag ---
// Order of params should be enforced by calling function
func RenameTag(srcName string, destName string) error {

	paramsMap := map[string]string{
		"s":    srcName,
		"dest": destName,
	}

	url := baseURL + "/rename-tag"
	if err := SetInoreader(url, paramsMap); err != nil {
		return errors.Wrapf(err, "Failed to rename %s tag to %s", srcName, destName)
	}

	return nil
}

// DeleteTag ---
func DeleteTag(tagName string) error {

	paramsMap := map[string]string{
		"s": tagName,
	}

	url := baseURL + "/disable-tag"
	if err := SetInoreader(url, paramsMap); err != nil {
		return errors.Wrapf(err, "Failed to delete tag %s", tagName)
	}

	return nil
}

// EditTag ---
// Order of params should be enforced by calling function
func EditTag(editTagParams []string) error {

	paramsMap := map[string]string{
		"a": editTagParams[0],
		"r": editTagParams[1],
		"i": editTagParams[2],
	}

	url := baseURL + "/edit-tag"
	if err := SetInoreader(url, paramsMap); err != nil {
		return errors.Wrap(err, "Failed to edit tag")
	}

	return nil
}

// MarkAllAsRead ---
// Order of params should be enforced by calling function
func MarkAllAsRead(ts string, streamid string) error {

	paramsMap := map[string]string{
		"ts": ts,
		"s":  streamid,
	}

	url := baseURL + "/mark-all-as-read"
	if err := SetInoreader(url, paramsMap); err != nil {
		return errors.Wrapf(err, "Failed to mark all items in %s as read", streamid)
	}

	return nil
}

// CreateActiveSearch ---
// Not implemented due to user account limitations

// DeleteActiveSearch ---
// Not implemented due to user account limitations
