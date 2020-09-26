package main

import (
	"github.com/hyperreal64/go-inoreader/oauthutil"
)

func main() {

	oauthutil.Init()

	// userInfo := &inoreader.UserInfo{}
	// err := inoreader.GetUserInfo(userInfo)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(userInfo.UserName)  // prints username
	// fmt.Println(userInfo.UserEmail) // prints user's signup email

	// unreadCounters := &inoreader.UnreadCounters{}
	// err := inoreader.GetUnreadCounters(unreadCounters)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, v := range unreadCounters.Unreadcounts {
	// 	fmt.Println(v)
	// }

	// subList := &inoreader.SubscriptionList{}
	// err := inoreader.GetSubscriptionList(subList)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, v := range subList.Subscriptions {
	// 	fmt.Println(v)
	// }

	// tagList := &inoreader.TagFolderList{}
	// err := inoreader.GetTagList(tagList)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// for _, v := range tagList.Tags {
	// 	fmt.Println(v)
	// }

	// streamContents := &inoreader.StreamContents{}
	// sid := "feed/http://feeds.arstechnica.com/arstechnica/science"
	// streamParams := []string{"3", "o", "", "", "", "", "", "false", sid}
	// err := inoreader.GetStreamContents(streamContents, streamParams)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// streamPrefsList := &inoreader.StreamPreferenceList{}
	// err := inoreader.GetStreamPrefsList(streamPrefsList)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// fmt.Println(streamPrefsList.Streamprefs)

	// editSubParams := []string{"edit", "feed/http://feeds.arstechnica.com/arstechnica/science", "", "user/-/label/Tech", ""}
	// if err := inoreader.EditSubscription(editSubParams); err != nil {
	// 	log.Fatalln(err)
	// }

	// editTagParams := []string{"", "user/-/state/com.google/read", "00000005fa676a18"}
	// if err := inoreader.EditTag(editTagParams); err != nil {
	// 	log.Fatalln(err)
	// }

	// timestamp := fmt.Sprint(time.Now().Unix())
	// if err := inoreader.MarkAllAsRead(timestamp, "feed/https://www.freecodecamp.org/news/rss/"); err != nil {
	// 	log.Fatalln(err)
	// }
}
