package main

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

// const (
// 	starredTag = "user/-/state/com.google/starred"
// 	savedTag   = "user/-/state/com.google/saved-web-pages"
// )

// StreamContents response
type StreamContents struct {
	Direction   string `json:"direction"`
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Updated     int    `json:"updated"`
	UpdatedUsec string `json:"updatedUsec"`
	Self        struct {
		Href string `json:"href"`
	} `json:"self"`
	Items        []Item `json:"items"`
	Continuation string `json:"continuation"`
}

// Item response
type Item struct {
	CrawlTimeMsec string   `json:"crawlTimeMsec"`
	TimestampUsec string   `json:"timestampUsec"`
	ID            string   `json:"id"`
	Categories    []string `json:"categories"`
	Title         string   `json:"title"`
	Published     int64    `json:"published"`
	Updated       int      `json:"updated"`
	Canonical     []struct {
		Href string `json:"href"`
	} `json:"canonical"`
	Author      string        `json:"author"`
	LikingUsers []interface{} `json:"likingUsers"`
	Comments    []interface{} `json:"comments"`
	CommentsNum int           `json:"commentsNum"`
	Annotations []interface{} `json:"annotations"`
	Origin      *Origin       `json:"origin"`
}

// Origin response
type Origin struct {
	StreamID string `json:"streamId"`
	Title    string `json:"title"`
	HTMLURL  string `json:"htmlUrl"`
}

// ItemIDs response
type ItemIDs struct {
	Items        []interface{} `json:"items"`
	ItemRefs     []interface{} `json:"itemRefs"`
	Continuation string        `json:"continuation"`
}

// ItemRefs response
type ItemRefs struct {
	ID              string        `json:"id"`
	DirectStreamIds []interface{} `json:"directStreamIds"`
	TimestampUsec   string        `json:"timestampUsec"`
}

// StreamPreferenceList response
type StreamPreferenceList struct {
	Streamprefs interface{} `json:"streamprefs"`
}

// Streamprefs response
type Streamprefs struct {
	UserStateComGoogleRoot []interface{} `json:"user/-/state/com.google/root"`
}

// UserStateComGoogleRoot response
type UserStateComGoogleRoot struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// GetStreamContents -- Gets stream contents based on set query parameters
func GetStreamContents(rc *resty.Client, params map[string]string) (sc *StreamContents, err error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Get(streamContentsURL)

	if err != nil {
		return nil, err
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), sc); err != nil {
		return nil, err
	}

	return sc, nil
}

// MarkAllAsRead -- Marks all items in stream as read; stream is specified in query parameters
func MarkAllAsRead(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(baseURL + "/mark-all-as-read")
	if err != nil {
		return err
	}

	return nil
}

// ExecMarkStreamRead -- Execute MarkAllAsRead for provided stream URL
func ExecMarkStreamAsRead(streamURL string) error {

	streamID := "feed/" + streamURL

	ctx, cancel := context.WithCancel(context.Background())
	rClient := oauth2RestyClient(ctx)
	defer cancel()

	params := map[string]string{
		"s": streamID,
	}

	if err := MarkAllAsRead(rClient, params); err != nil {
		return errors.Wrapf(err, "Could not mark stream as read: %s", streamID)
	}

	return nil
}
