package main

import "encoding/json"

// UserInfo ---
type UserInfo struct {
	UserID              string `json:"userId"`
	UserName            string `json:"userName"`
	UserProfileID       string `json:"userProfileId"`
	UserEmail           string `json:"userEmail"`
	IsBloggerUser       bool   `json:"isBloggerUser"`
	SignupTimeSec       int64  `json:"signupTimeSec"`
	IsMultiLoginEnabled bool   `json:"isMultiLoginEnabled"`
}

// QuickAdd --- JSON object returned from POST request
type QuickAdd struct {
	Query      string `json:"query"`
	NumResults int    `json:"numResults"`
	StreamID   string `json:"streamId"`
	StreamName string `json:"streamName"`
}

// UnreadCounters ---
type UnreadCounters struct {
	Max          int `json:"max"`
	Unreadcounts []struct {
		ID                      string      `json:"id"`
		Count                   json.Number `json:"count"`
		NewestItemTimestampUsec string      `json:"newestItemTimestampUsec"`
	} `json:"unreadcounts"`
}

// SubscriptionList ---
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

// TagFolderList ---
type TagFolderList struct {
	Tags []struct {
		ID          string `json:"id"`
		Sortid      string `json:"sortid"`
		UnreadCount int64  `json:"unread_count"`
		Type        string `json:"type"`
	} `json:"tags"`
}

// StreamContents ---
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

// Item --- of StreamContents
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

// Origin --- of Items of StreamContents
type Origin struct {
	StreamID string `json:"streamId"`
	Title    string `json:"title"`
	HTMLURL  string `json:"htmlUrl"`
}

// ItemIDs ---
type ItemIDs struct {
	Items        []interface{} `json:"items"`
	ItemRefs     []interface{} `json:"itemRefs"`
	Continuation string        `json:"continuation"`
}

// ItemRefs --- of ItemIDs
type ItemRefs struct {
	ID              string        `json:"id"`
	DirectStreamIds []interface{} `json:"directStreamIds"`
	TimestampUsec   string        `json:"timestampUsec"`
}

// StreamPreferenceList ---
type StreamPreferenceList struct {
	Streamprefs interface{} `json:"streamprefs"`
}

// Streamprefs --- of StreamPreferenceList
type Streamprefs struct {
	UserStateComGoogleRoot []interface{} `json:"user/-/state/com.google/root"`
}

// UserStateComGoogleRoot --- of StreamPrefs of StreamPreferenceList
type UserStateComGoogleRoot struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
