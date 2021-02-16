package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func printTagsFolders(tagList *TagFolderList, onlyUnread bool, option string) {

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
