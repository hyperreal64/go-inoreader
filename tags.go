package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// FIXME: Refactor
// See about making the four functions below into a single function.
// Use maybe a strucc or map[string]bool to check if the corresponding
// argument was passed to the calling function.
func printTagFolderList(tagList *TagFolderList) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Label", "Type", "# Unread"})

	for _, v := range tagList.Tags {
		label := strings.Split(v.ID, "/")
		labelSuffix := label[len(label)-1]
		table.Append([]string{labelSuffix, v.Type, strconv.FormatInt(v.UnreadCount, 10)})
	}
	table.Render()
}

// Print only folders/tags with unread counts
func printTagFolderListUnread(tagList *TagFolderList) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Label", "Type", "# Unread"})

	for _, v := range tagList.Tags {
		count := v.UnreadCount
		label := strings.Split(v.ID, "/")
		labelSuffix := label[len(label)-1]
		if count > 0 {
			table.Append([]string{labelSuffix, v.Type, strconv.FormatInt(v.UnreadCount, 10)})
		}
	}
	table.Render()
}

// Print tags XOR folders
func printTagsXorFolders(tagList *TagFolderList, option string) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{option})

	switch option {
	case "tags":
		for _, v := range tagList.Tags {
			if v.Type == "tag" {
				label := strings.Split(v.ID, "/")
				labelSuffix := label[len(label)-1]
				table.Append([]string{labelSuffix})
			}
		}

	case "folders":
		for _, v := range tagList.Tags {
			if v.Type == "folder" {
				label := strings.Split(v.ID, "/")
				labelSuffix := label[len(label)-1]
				table.Append([]string{labelSuffix})
			}
		}
	}
	table.Render()
}

// Print tags XOR folders with unread counts
func printTagsXorFoldersUnread(tagList *TagFolderList, option string) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{option, "# unread"})

	switch option {
	case "tags":
		for _, v := range tagList.Tags {
			if v.Type == "tag" && v.UnreadCount > 0 {
				label := strings.Split(v.ID, "/")
				labelSuffix := label[len(label)-1]
				table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
			}
		}

	case "folders":
		for _, v := range tagList.Tags {
			if v.Type == "folder" && v.UnreadCount > 0 {
				label := strings.Split(v.ID, "/")
				labelSuffix := label[len(label)-1]
				table.Append([]string{labelSuffix, strconv.FormatInt(v.UnreadCount, 10)})
			}
		}
	}
	table.Render()
}
