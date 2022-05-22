package subscription

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// API URLs for resty.Client requests
const (
	addSubURL         = "https://www.inoreader.com/reader/api/0/subscription/quickadd"
	editSubURL        = "https://www.inoreader.com/reader/api/0/subscription/edit"
	unreadCountersURL = "https://www.inoreader.com/reader/api/0/unread-count"
	subListURL        = "https://www.inoreader.com/reader/api/0/subscription/list"
)

// QuickAdd JSON response
type QuickAdd struct {
	Query      string `json:"query"`
	NumResults int    `json:"numResults"`
	StreamID   string `json:"streamId"`
	StreamName string `json:"streamName"`
}

// UnreadCounters JSON response
type UnreadCounters struct {
	Max          int `json:"max"`
	Unreadcounts []struct {
		ID                      string      `json:"id"`
		Count                   json.Number `json:"count"`
		NewestItemTimestampUsec string      `json:"newestItemTimestampUsec"`
	} `json:"unreadcounts"`
}

// SubscriptionList JSON response
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

// Quick add a subscription as specified in the query parameters.
// Unlike other POST requests to the Inoreader API server, this one returns
// a JSON response, which gets stored into a QuickAdd struct.
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

// Edit subscription specified in query parameters. Sends a POST request.
func EditSubscription(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(editSubURL)
	if err != nil {
		return err
	}

	return nil
}

// Get list of subscriptions. Sends a GET request and returns JSON response as
// SubscriptionList struct.
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

// Get the number of unread items for subscriptions. Sends a GET request and
// returns JSON response as UnreadCounters struct.
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
