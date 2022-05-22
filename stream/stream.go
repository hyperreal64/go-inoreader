package stream

import (
	"github.com/go-resty/resty/v2"
)

// API URLs for resty.Client requests
const (
	streamContentsURL = "https://www.inoreader.com/reader/api/0/stream/contents"
	itemIDsURL        = "https://www.inoreader.com/reader/api/0/stream/items/ids"
	streamPrefsURL    = "https://www.inoreader.com/reader/api/0/preference/stream/list"
	streamPrefsSetURL = "https://www.inoreader.com/reader/api/0/preference/stream/set"
	markAllReadURL    = "https://www.inoreader.com/reader/api/0/mark-all-as-read"
)

// StreamContents JSON response
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

// Item JSON response
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

// Origin JSON response
type Origin struct {
	StreamID string `json:"streamId"`
	Title    string `json:"title"`
	HTMLURL  string `json:"htmlUrl"`
}

// ItemIDs JSON response
type ItemIDs struct {
	Items        []interface{} `json:"items"`
	ItemRefs     []interface{} `json:"itemRefs"`
	Continuation string        `json:"continuation"`
}

// ItemRefs JSON response
type ItemRefs struct {
	ID              string        `json:"id"`
	DirectStreamIds []interface{} `json:"directStreamIds"`
	TimestampUsec   string        `json:"timestampUsec"`
}

// StreamPreferenceList JSON response
type StreamPreferenceList struct {
	Streamprefs interface{} `json:"streamprefs"`
}

// Streamprefs JSON response
type Streamprefs struct {
	UserStateComGoogleRoot []interface{} `json:"user/-/state/com.google/root"`
}

// UserStateComGoogleRoot JSON response
type UserStateComGoogleRoot struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

// Gets stream contents based on set query parameters and returns a struct
// containing StreamContents JSON response
func GetStreamContents(rc *resty.Client, params map[string]string) (sc *StreamContents, err error) {

	resp, err := rc.R().
		SetQueryParams(params).
		Get(streamContentsURL)

	if err != nil {
		return nil, err
	}

	if err := resty.Unmarshalc(rc, "application/json", resp.Body(), &sc); err != nil {
		return nil, err
	}

	return sc, nil
}

// Marks all items in stream as read; stream is specified in query parameters
func MarkAllAsRead(rc *resty.Client, params map[string]string) error {

	_, err := rc.R().
		SetQueryParams(params).
		Post(markAllReadURL)
	if err != nil {
		return err
	}

	return nil
}
