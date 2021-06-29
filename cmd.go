package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

// UserInfoCmd ---
type UserInfoCmd struct{}

// Run (*UserInfoCmd) ---
func (u *UserInfoCmd) Run(ctx *kong.Context) error {

	if err := printUserInfo(); err != nil {
		return err
	}

	return nil
}

// LoginCmd ---
type LoginCmd struct{}

// Run (*LoginCmd) ---
func (l *LoginCmd) Run(ctx *kong.Context) error {

	Init()

	return nil
}

// ListCmd ---
type ListCmd struct {
	Subscriptions SubListCmd      `cmd:"" aliases:"subs" help:"List subscriptions"`
	Tags          TagListCmd      `cmd:"" help:"List tags and/or folders"`
	Stream        StreamListCmd   `cmd:"" help:"List stream contents"`
	Starred       StarredListCmd  `cmd:"" aliases:"star" help:"List starred articles"`
	WebPages      WebPagesListCmd `cmd:"" aliases:"wp" help:"List web pages saved on Inoreader"`
}

// Run (*ListCmd) ---
func (l *ListCmd) Run(ctx *kong.Context) error {

	return nil
}

// SubListCmd ---
type SubListCmd struct {
	All    bool `required,name:"all" short:"a" xor:"list" help:"List all subscriptions"`
	Unread bool `required,name:"unread" short:"u" xor:"list" help:"List only unread subscriptions"`
}

// Run (*SubListCmd) ---
func (s *SubListCmd) Run(ctx *kong.Context) error {

	switch true {
	case s.All:
		if err := printSubList(false); err != nil {
			return err
		}

	case s.Unread:
		if err := printSubList(true); err != nil {
			return err
		}

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// TagListCmd ---
type TagListCmd struct {
	All    bool   `required,name:"all" short:"a" xor:"list" help:"List all tags"`
	Unread bool   `required,name:"unread" short:"u" xor:"list" help:"List only unread tags"`
	Type   string `arg,optional,short:"t" placeholder:"tags|folders" help:"List either only tags or only folders"`
}

// Run (*TagListCmd) ---
func (t *TagListCmd) Run(ctx *kong.Context) error {

	switch true {
	case t.All:
		if err := printTagsFolders(false, t.Type); err != nil {
			return err
		}

	case t.Unread:
		if err := printTagsFolders(true, t.Type); err != nil {
			return err
		}

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// StreamListCmd ---
type StreamListCmd struct {
	URL           string `arg:"" required,name:"url" help:"Specify URL of feed"`
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	URLS          bool   `arg,optional,name:"urls" short:"u" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" xor:"stream" help:"Display items with their IDs"`
	Dates         bool   `arg,optional,name:"dates" short:"d" xor:"stream" help:"Display items with their timestamps"`
}

// Run (*StreamListCmd) ---
func (s *StreamListCmd) Run(ctx *kong.Context) error {

	switch true {
	case s.URLS:
		if err := printStreamContentsWithURL(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, s.URL); err != nil {
			return err
		}

	case s.IDS:
		if err := printStreamContentsWithIDs(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, s.URL); err != nil {
			return err
		}

	case s.Dates:
		if err := printStreamContentsWithDate(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, s.URL); err != nil {
			return err
		}

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// StarredListCmd ---
type StarredListCmd struct {
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	URLS          bool   `arg,optional,name:"urls" short:"u" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" xor:"stream" help:"Display items with their IDs"`
	Dates         bool   `arg,optional,name:"dates" short:"d" xor:"stream" help:"Display items with their timestamps"`
}

// Run (*StarredListCmd) ---
func (s *StarredListCmd) Run(ctx *kong.Context) error {

	switch true {
	case s.URLS:
		err := printStreamContentsWithURL(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}

	case s.IDS:
		err := printStreamContentsWithIDs(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}

	case s.Dates:
		err := printStreamContentsWithDate(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// WebPagesListCmd ---
type WebPagesListCmd struct {
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	URLS          bool   `arg,optional,name:"urls" short:"u" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" xor:"stream" help:"Display items with their IDs"`
	Dates         bool   `arg,optional,name:"dates" short:"d" xor:"stream" help:"Display items with their timestamps"`
}

// Run (*WebPagesListCmd) ---
func (w *WebPagesListCmd) Run(ctx *kong.Context) error {

	switch true {
	case w.URLS:
		err := printStreamContentsWithURL(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
			return err
		}

	case w.IDS:
		err := printStreamContentsWithIDs(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
			return err
		}

	case w.Dates:
		err := printStreamContentsWithDate(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
			return err
		}

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// SubscriptionCmd ---
type SubscriptionCmd struct {
	Add           AddSubCmd        `cmd:"" help:"Add a subscription"`
	Unsubscribe   UnsubscribeCmd   `cmd:"" aliases:"un" help:"Unsubscribe from <url>"`
	SetTitle      SetTitleCmd      `cmd:"" aliases:"st" help:"Change the title of a subscription"`
	AddToFolder   AddToFolderCmd   `cmd:"" aliases:"af" help:"Add subscription to folder"`
	RemFromFolder RemFromFolderCmd `cmd:"" aliases:"rf" help:"Remove subscription from folder"`
}

// Run (*SubscriptionCmd) ---
func (s *SubscriptionCmd) Run(ctx *kong.Context) error {

	return nil
}

// AddSubCmd ---
type AddSubCmd struct {
	URL string `arg:"" required,help:"URL of subscription feed"`
}

// Run (*AddSubCmd) ---
func (a *AddSubCmd) Run(ctx *kong.Context) error {

	if err := execAddSub(a.URL); err != nil || a.URL == "" {
		return err
	}

	fmt.Printf("%s added to Inoreader\n", a.URL)

	return nil
}

// UnsubscribeCmd ---
type UnsubscribeCmd struct {
	URL string `arg,required,help:"Unsubscribe from <url>"`
}

// Run (*UnsubscribeCmd) ---
func (u *UnsubscribeCmd) Run(ctx *kong.Context) error {

	if err := execUnsubscribe(u.URL); err != nil || u.URL == "" {
		return err
	}

	fmt.Printf("%s removed from Inoreader\n", u.URL)

	return nil
}

// SetTitleCmd ---
type SetTitleCmd struct {
	Title string `arg,required,help:"New title of subscription"`
	URL   string `arg,required,help:"URL of subscription to change"`
}

// Run (*SetTitleCmd) ---
func (t *SetTitleCmd) Run(ctx *kong.Context) error {

	if err := execSetSubTitle(t.Title, t.URL); err != nil || t.Title == "" || t.URL == "" {
		return err
	}

	fmt.Printf("Changed title of feed to: %s\n", t.Title)

	return nil
}

// AddToFolderCmd ---
type AddToFolderCmd struct {
	Folder string `arg,required,help:"Folder to add subscription to"`
	URL    string `arg,required,help:"URL of subscription"`
}

// Run (*AddToFolderCmd) ---
func (a *AddToFolderCmd) Run(ctx *kong.Context) error {

	if err := execAddSubToFolder(a.Folder, a.URL); err != nil || a.Folder == "" || a.URL == "" {
		return err
	}

	fmt.Printf("+ Added %s to %s\n", a.URL, a.Folder)

	return nil
}

// RemFromFolderCmd ---
type RemFromFolderCmd struct {
	Folder string `arg,required,short:"f" help:"Folder to remove subscription from"`
	URL    string `arg,required,short:"u" help:"URL of subscription"`
}

// Run (*RemFromFolderCmd) ---
func (r *RemFromFolderCmd) Run(ctx *kong.Context) error {

	if err := execRemSubFromFolder(r.Folder, r.URL); err != nil {
		return err
	}

	fmt.Printf("- Removed %s from %s\n", r.URL, r.Folder)

	return nil
}

// TagsCmd ---
type TagsCmd struct {
	RenameTag RenameTagCmd `cmd:"" aliases:"mv" help:"Rename a tag"`
	DeleteTag DeleteTagCmd `cmd:"" aliases:"rm" help:"Delete a tag"`
}

// Run (*TagsCmd) ---
func (t *TagsCmd) Run(ctx *kong.Context) error {

	return nil
}

// RenameTagCmd ---
type RenameTagCmd struct {
	Src  string `arg,required,short:"s" help:"Source tag name"`
	Dest string `arg,required,short:"d" help:"Destination tag name"`
}

// Run (*RenameTagCmd) ---
func (r *RenameTagCmd) Run(ctx *kong.Context) error {

	if err := execRenameTag(r.Src, r.Dest); err != nil {
		return err
	}

	fmt.Printf("Renamed %s tag to %s\n", r.Src, r.Dest)

	return nil
}

// DeleteTagCmd ---
type DeleteTagCmd struct {
	Delete string `arg,required,help:"Tag name to delete"`
}

// Run (*DeleteTagCmd) ---
func (d *DeleteTagCmd) Run(ctx *kong.Context) error {

	if err := execDelTag(d.Delete); err != nil {
		return err
	}

	fmt.Printf("Deleted %s tag\n", d.Delete)

	return nil
}

// MarkItemCmd ---
type MarkItemCmd struct {
	Read   bool   `arg,xor:"read" help:"Mark item as read"`
	Unread bool   `arg,xor:"read" help:"Mark item as unread"`
	Star   bool   `arg,xor:"star" help:"Star an article"`
	Unstar bool   `arg,xor:"star" help:"Unstar an article"`
	ItemID string `arg,required,help:"Item ID"`
}

// Run (*MarkItemCmd) ---
func (m *MarkItemCmd) Run(ctx *kong.Context) error {

	switch true {
	case m.Read:
		if err := execEditTagRead(m.ItemID, true); err != nil {
			return err
		}

		fmt.Printf("Marked %s as read\n", m.ItemID)

	case m.Unread:
		if err := execEditTagRead(m.ItemID, false); err != nil {
			return err
		}

		fmt.Printf("Attempted to mark %s as unread.\n", m.ItemID)
		fmt.Println("Please note that if the item's timestamp is older than the first unread timestamp of its feed,")
		fmt.Println("then it cannot be marked unread.")

	case m.Star:
		if err := execEditTagStar(m.ItemID, true); err != nil {
			return err
		}

		fmt.Printf("Marked %s as starred.\n", m.ItemID)

	case m.Unstar:
		if err := execEditTagStar(m.ItemID, false); err != nil {
			return err
		}

		fmt.Printf("Unstarred %s\n", m.ItemID)

	default:
		if err := ctx.PrintUsage(true); err != nil {
			return err
		}
	}

	return nil
}

// MarkStreamReadCmd ---
type MarkStreamReadCmd struct {
	StreamID string `arg,required,short:"s" help:"Enter stream ID or feed URL to mark as read"`
}

// Run (*MarkStreamReadCmd) ---
func (r *MarkStreamReadCmd) Run(ctx *kong.Context) error {

	if err := execMarkStreamAsRead(r.StreamID); err != nil {
		return err
	}

	fmt.Printf("Marked %s as read\n", r.StreamID)

	return nil
}

// ExamplesCmd ---
type ExamplesCmd struct{}

// Run (*ExamplesCmd) ---
func (e *ExamplesCmd) Run(ctx *kong.Context) error {

	printCmdExamples()

	return nil
}

var cli struct {
	Examples       ExamplesCmd       `cmd:"" aliases:"ex" help:"Print examples of go-inoreader commands"`
	List           ListCmd           `cmd:"" aliases:"ls" help:"List Inoreader subscriptions, tags/folders, stream contents, starred articles, or saved web pages"`
	Login          LoginCmd          `cmd:"" help:"Login and initiate Oauth flow"`
	MarkItem       MarkItemCmd       `cmd:"" aliases:"mark" help:"Mark items as read/unread or starred/unstarred"`
	MarkStreamRead MarkStreamReadCmd `cmd:"" help:"Mark all items in stream as read"`
	Subscription   SubscriptionCmd   `cmd:"" aliases:"sub" help:"Query subscriptions"`
	Tags           TagsCmd           `cmd:"" help:"Query tags and folders"`
	UserInfo       UserInfoCmd       `cmd:"" help:"Print Inoreader user information"`
}

func main() {

	ctx := kong.Parse(&cli,
		kong.Name("go-inoreader"),
		kong.Description("An Inoreader API client for GoLang"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary: false,
			Compact:      true,
			Tree:         false,
		}))
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}
