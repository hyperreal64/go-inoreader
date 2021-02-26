package main

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
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

	getRestyErr string = "Could not get resty client"
)

// GetUserInfo ---
func getUserInfo(rc *resty.Client) (*UserInfo, error) {

	resp, err := rc.R().Get(userInfoURL)
	if err != nil {
		return nil, errors.Wrap(err, "Could not get user info")
	}

	userInfo := &UserInfo{}
	if err := json.Unmarshal(resp.Body(), userInfo); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", userInfo)
	}

	return userInfo, nil
}

// QuickAddSubscription ---
func quickAddSubscription(rc *resty.Client, params map[string]string) (*QuickAdd, error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Post(addSubURL)
	if err != nil {
		return nil, err
	}

	quickAdd := &QuickAdd{}
	if err := json.Unmarshal(resp.Body(), quickAdd); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", quickAdd)
	}

	return quickAdd, nil
}

// EditSubscription ---
func editSubscription(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(editSubURL)
	if err != nil {
		return err
	}

	return nil
}

func getSubscriptionList(rc *resty.Client) (*SubscriptionList, error) {

	resp, err := rc.R().Get(subListURL)
	if err != nil {
		return nil, err
	}

	subList := &SubscriptionList{}
	if err := json.Unmarshal(resp.Body(), &subList); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", subList)
	}

	return subList, nil
}

func getUnreadCounters(rc *resty.Client) (*UnreadCounters, error) {

	resp, err := rc.R().Get(unreadCountersURL)
	if err != nil {
		return nil, err
	}

	unreadCounters := &UnreadCounters{}
	if err = json.Unmarshal(resp.Body(), &unreadCounters); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object %v", unreadCounters)
	}

	return unreadCounters, nil
}

func getTagList(rc *resty.Client) (*TagFolderList, error) {

	resp, err := rc.R().Get(tagListURL)
	if err != nil {
		return nil, err
	}

	tagList := &TagFolderList{}
	if err := json.Unmarshal(resp.Body(), tagList); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", tagList)
	}

	return tagList, nil
}

func getStreamContents(rc *resty.Client, params map[string]string) (*StreamContents, error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Get(streamContentsURL)
	if err != nil {
		return nil, err
	}

	streamContents := &StreamContents{}
	if err := json.Unmarshal(resp.Body(), streamContents); err != nil {
		return nil, errors.Wrapf(err, "Could not unmarshal JSON object: %v", streamContents)
	}

	return streamContents, nil
}

// GetItemIDs ---
func getItemIDs(rc *resty.Client, itemIDs *ItemIDs, params map[string]string) error {

	resp, err := rc.R().
		SetQueryParams(params).
		Get(itemIDsURL)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), itemIDs); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", itemIDs)
	}

	return nil
}

// GetStreamPrefsList ---
func getStreamPrefsList(rc *resty.Client, streamPrefsList *StreamPreferenceList) error {

	resp, err := rc.R().Get(streamPrefsURL)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(resp.Body(), streamPrefsList); err != nil {
		return errors.Wrapf(err, "Could not unmarshal JSON object: %v", streamPrefsList)
	}

	return nil
}

// SetStreamPrefs ---
func setStreamPrefs(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(streamPrefsSetURL)
	if err != nil {
		return err
	}

	return nil
}

// RenameTag ---
func renameTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/rename-tag")
	if err != nil {
		return err
	}

	return nil
}

// DeleteTag ---
func deleteTag(rc *resty.Client, tagName string) error {

	_, err := rc.R().
		SetQueryParams(map[string]string{
			"s": tagName,
		}).
		Post(baseURL + "/disable-tag")
	if err != nil {
		return err
	}

	return nil
}

// EditTag ---
func editTag(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/edit-tag")
	if err != nil {
		return err
	}

	return nil
}

// MarkAllAsRead ---
func markAllAsRead(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/mark-all-as-read")
	if err != nil {
		return err
	}

	return nil
}

// CreateActiveSearch ---
// Not implemented due to user account limitations

// DeleteActiveSearch ---
// Not implemented due to user account limitations
