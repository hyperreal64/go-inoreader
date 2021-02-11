package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// TODO: printTagFolderList with option for tag, folder, or both
//   Command: inoreader tags -o <tag|folder|both>
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
