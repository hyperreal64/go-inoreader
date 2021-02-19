package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func printTagsFolders(onlyUnread bool, option string) {

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	tagList, err := getTagList(rClient)
	if err != nil {
		log.Fatalln(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{option, "# unread"})

	switch option {
	case "tags":
		for _, v := range tagList.Tags {
			if onlyUnread == true {
				if v.Type == "tag" && v.UnreadCount > 0 {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			} else {
				if v.Type == "tag" {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			}
		}

	case "folders":
		for _, v := range tagList.Tags {
			if onlyUnread == true {
				if v.Type == "folder" && v.UnreadCount > 0 {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			} else {
				if v.Type == "folder" {
					label := strings.Split(v.ID, "/")
					labelSuffix := label[len(label)-1]
					table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
				}
			}
		}
	}
	table.Render()
}

// TODO: getItemID()
func execEditTag(addTag string, remTag string, itemID string) string {

	params := map[string]string{
		"a": addTag,
		"r": remTag,
		"i": itemID,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := editTag(rClient, params); err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("Successfully edited tags for item: %s\n", itemID)
}

func execRenameTag(src string, dest string) string {

	params := map[string]string{
		"s":    src,
		"dest": dest,
	}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := renameTag(rClient, params); err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("Successfully renamed tag from %s to %s\n", src, dest)
}

func execDelTag(tagName string) string {

	params := map[string]string{"s": tagName}

	rClient, err := config2Client()
	if err != nil {
		log.Fatalln(err)
	}

	if err := renameTag(rClient, params); err != nil {
		log.Fatalln(err)
	}

	return fmt.Sprintf("Successfully deleted tag: %s\n", tagName)
}
