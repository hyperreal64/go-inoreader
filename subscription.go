package main

import (
	"context"
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// QuickAdd response
type QuickAdd struct {
	Query      string `json:"query"`
	NumResults int    `json:"numResults"`
	StreamID   string `json:"streamId"`
	StreamName string `json:"streamName"`
}

// UnreadCounters response
type UnreadCounters struct {
	Max          int `json:"max"`
	Unreadcounts []struct {
		ID                      string      `json:"id"`
		Count                   json.Number `json:"count"`
		NewestItemTimestampUsec string      `json:"newestItemTimestampUsec"`
	} `json:"unreadcounts"`
}

// SubscriptionList response
type SubscriptionList struct {
	Subscriptions []struct {
		ID            string        `json:"id"`
		FeedType      string        `json:"feedType"`
		Title         string        `json:"title"`
		Categories    []interface{} `json:"categories"`
		Sortid        string        `json:"sortid"`
		Firstitemmsec int64         `json:"firstitemmsec"`
		URL           string        `json:"url"`
		HTMLURL       string        `json:"htmlUrl"`
		IconURL       string        `json:"iconUrl"`
	} `json:"subscriptions"`
}

// QuickAddSubscription --- Quick add a subscription as specified in query parameters.
// Unlike other POST requests to the Inoreader API server, this one returns a JSON response.
// Parameters:
// rc --> resty.Client
// params --> query parameters that contain the subscription's URL
// Returns: QuickAdd as JSON object, or error
func QuickAddSubscription(rc *resty.Client, params map[string]string) (quickadd *QuickAdd, err error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Post(addSubURL)
	if err != nil {
		return nil, err
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &quickadd); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object: %v", quickadd)
	}

	return quickadd, nil
}

// EditSubscription -- Edit subscription specified in query parameters.
// Parameters:
// rc --> resty.Client
// params --> query parameters that contain subscription URL and an action to take
// Returns: error on error
func EditSubscription(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(editSubURL)
	if err != nil {
		return err
	}

	return nil
}

// GetSubscriptionList -- Get list of subscriptions.
func GetSubscriptionList(rc *resty.Client) (sublist *SubscriptionList, err error) {

	resp, err := rc.R().Get(subListURL)
	if err != nil {
		return nil, err
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &sublist); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object: %v", sublist)
	}

	return sublist, nil
}

// GetUnreadCounters -- Get the number of unread items for subscriptions.
func GetUnreadCounters(rc *resty.Client) (uc *UnreadCounters, err error) {

	resp, err := rc.R().Get(unreadCountersURL)
	if err != nil {
		return nil, err
	}

	if err = resty.Unmarshalc(rc, "application/json", resp.Body(), &uc); err != nil {
		return nil, errors.Wrapf(err, "Unable to unmarshal JSON object %v", uc)
	}

	return uc, nil
}

// ExecAddSub -- Execute quick add subscription
func ExecAddSub(url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"quickadd": streamID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	quickAdd, err := QuickAddSubscription(rc, params)
	if err != nil {
		return err
	}

	if quickAdd.NumResults == 0 {
		return errors.New("Please check if the URL is correct")
	}

	return nil
}

// ExecUnsubscribe -- Execute unsubscribe for subscription at URL
func ExecUnsubscribe(url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "unsubscribe",
		"s":  streamID,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditSubscription(rc, params); err != nil {
		return errors.Wrapf(err, "Unable to unsubscribe from subscription %s", url)
	}

	return nil
}

// ExecSetSubTitle -- Execute EditSubcription to set the subscription's title
func ExecSetSubTitle(title string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"t":  title,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditSubscription(rc, params); err != nil {
		return errors.Wrapf(err, "Unable to set title %s on subscription %s", title, url)
	}

	return nil
}

// ExecAddSubToFolder -- Execute EditSubscription to add a subscription to specified folder
func ExecAddSubToFolder(folder string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"a":  folder,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditSubscription(rc, params); err != nil {
		return errors.Wrapf(err, "Unable to add subscription %s to folder %s", url, folder)
	}

	return nil
}

// ExecRemSubFromFolder --- Execute EditSubscription to remove subscription from folder
func ExecRemSubFromFolder(folder string, url string) error {

	streamID := "feed/" + url
	params := map[string]string{
		"ac": "edit",
		"s":  streamID,
		"r":  folder,
	}

	ctx, cancel := context.WithCancel(context.Background())
	rc := oauth2RestyClient(ctx)
	defer cancel()

	if err := EditSubscription(rc, params); err != nil {
		return errors.Wrapf(err, "Unable to remove subscription %s from folder %s", url, folder)
	}

	return nil
}
