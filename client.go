package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// GetSubscriptionTitle ---
func GetSubscriptionTitle(sid string, client *http.Client) string {

	subscriptionList := &SubscriptionList{}
	if err := GetSubscriptionList(client, subscriptionList); err != nil {
		log.Fatalln(err)
	}

	var title string

	for _, v := range subscriptionList.Subscriptions {
		if v.ID == sid {
			title = v.Title
		}
	}

	return title
}

// ListUnreadCounters ---
func ListUnreadCounters(client *http.Client) {

	unreadCounters := &UnreadCounters{}
	if err := GetUnreadCounters(client, unreadCounters); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Unread items")
	fmt.Println(strings.Repeat("-", 30))
	for _, v := range unreadCounters.Unreadcounts {
		count, err := v.Count.Int64()
		if err != nil {
			log.Fatalln(err)
		}

		statePrefix := "user/1005869311/state/com.google/"
		if count > 0 && v.ID != statePrefix+"reading-list" && v.ID != statePrefix+"starred" {
			title := GetSubscriptionTitle(v.ID, client)
			fmt.Printf("%-5d %s\n", count, title)
		}
	}
}
