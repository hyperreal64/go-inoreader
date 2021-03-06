package main

import (
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

// SubscriptionCmd ---
type SubscriptionCmd struct {
	Add           AddSubCmd        `cmd:"" help:"Add a subscription"`
	List          ListSubsCmd      `cmd:"" aliases:"ls" help:"List subscriptions"`
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

	if err := execAddSub(a.URL); err != nil {
		return err
	}

	return nil
}

// ListSubsCmd ---
type ListSubsCmd struct {
	All    bool `required,name:"all" short:"a" xor:"list" help:"List all subscriptions"`
	Unread bool `required,name:"unread" short:"u" xor:"list" help:"List only unread subscriptions"`
}

// Run (*ListCmd) ---
func (l *ListSubsCmd) Run(ctx *kong.Context) error {

	if l.All {
		if err := printSubList(false); err != nil {
			return err
		}
	}

	if l.Unread {
		if err := printSubList(true); err != nil {
			return err
		}
	}

	if !l.All && !l.Unread {
		ctx.PrintUsage(true)
	}

	return nil
}

// UnsubscribeCmd ---
type UnsubscribeCmd struct {
	URL string `arg,required,help:"Unsubscribe from <url>"`
}

// Run (*UnsubscribeCmd) ---
func (u *UnsubscribeCmd) Run(ctx *kong.Context) error {

	if err := execUnsubscribe(u.URL); err != nil {
		return err
	}

	return nil
}

// SetTitleCmd ---
type SetTitleCmd struct {
	Title string `arg,required,help:"New title of subscription"`
	URL   string `arg,required,help:"URL of subscription to change"`
}

// Run (*SetTitleCmd) ---
func (t *SetTitleCmd) Run(ctx *kong.Context) error {

	if err := execSetSubTitle(t.Title, t.URL); err != nil {
		return err
	}

	return nil
}

// AddToFolderCmd ---
type AddToFolderCmd struct {
	Folder string `arg,required,help:"Folder to add subscription to"`
	URL    string `arg,required,help:"URL of subscription"`
}

// Run (*AddToFolderCmd) ---
func (a *AddToFolderCmd) Run(ctx *kong.Context) error {

	if err := execAddSubToFolder(a.Folder, a.URL); err != nil {
		return err
	}

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

	return nil
}

// TagsCmd ---
type TagsCmd struct {
	List      ListTagsCmd  `cmd:"" aliases:"ls" help:"List tags and/or folders"`
	RenameTag RenameTagCmd `cmd:"" aliases:"mv" help:"Rename a tag"`
	DeleteTag DeleteTagCmd `cmd:"" aliases:"rm" help:"Delete a tag"`
}

// Run (*TagsCmd) ---
func (t *TagsCmd) Run(ctx *kong.Context) error {

	return nil
}

// ListTagsCmd ---
type ListTagsCmd struct {
	All    bool   `required,name:"all" short:"a" xor:"list" help:"List all tags"`
	Unread bool   `required,name:"unread" short:"u" xor:"list" help:"List only unread tags"`
	Type   string `arg,optional,short:"t" placeholder:"tags|folders" help:"List either only tags or only folders"`
}

// Run (*ListTagsCmd) ---
func (t *ListTagsCmd) Run(ctx *kong.Context) error {

	if t.All {
		if err := printTagsFolders(false, t.Type); err != nil {
			return err
		}
	}

	if t.Unread {
		if err := printTagsFolders(true, t.Type); err != nil {
			return err
		}
	}

	if !t.All && !t.Unread {
		ctx.PrintUsage(true)
	}

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

	return nil
}

// MarkItemReadCmd ---
type MarkItemReadCmd struct {
	Read   bool   `arg,xor:"read" help:"Mark item as read"`
	Unread bool   `arg,xor:"read" help:"Mark item as unread"`
	Star   bool   `arg,xor:"starred" help:"Add article to Starred"`
	Unstar bool   `arg,xor:"starred" help:"Remove article from Starred"`
	ItemID string `arg,required,help:"Item ID"`
}

// Run (*MarkItemCmd) ---
func (m *MarkItemReadCmd) Run(ctx *kong.Context) error {

	if m.Read {
		if err := execEditTagRead(m.ItemID, true); err != nil {
			return err
		}
	}

	if m.Unread {
		if err := execEditTagRead(m.ItemID, false); err != nil {
			return err
		}
	}

	if m.Star {
		if err := execEditTagStar(m.ItemID, true); err != nil {
			return err
		}
	}

	if m.Unstar {
		if err := execEditTagStar(m.ItemID, false); err != nil {
			return err
		}
	}

	if !m.Read && !m.Unread && !m.Star && !m.Unstar {
		ctx.PrintUsage(true)
	}

	return nil
}

// StreamCmd ---
type StreamCmd struct {
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	StreamID      string `arg,required,name:"stream-id" short:"s" help:"Specify stream ID of feed"`
	URLS          bool   `arg,optional,name:"urls" short:"u" default:"false" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" default:"false" xor:"stream" help:"Display items with their IDs"`
}

// Run (*StreamCmd) ---
func (sc *StreamCmd) Run(ctx *kong.Context) error {
	if sc.URLS {
		if err := printStreamContentsWithURL(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID); err != nil {
			return err
		}
	} else if sc.IDS {
		if err := printStreamContentsWithIDs(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID); err != nil {
			return err
		}
	} else {
		if err := printStreamContentsWithDate(sc.Num, sc.Order, sc.ExcludeTarget, sc.IncludeTarget, sc.StreamID); err != nil {
			return err
		}
	}

	return nil
}

// StarredCmd ---
type StarredCmd struct {
	List ListStarredCmd `cmd:"" aliases:"ls" help:"List starred articles"`
}

// Run (*StarredCmd) ---
func (s *StarredCmd) Run(ctx *kong.Context) error {

	return nil
}

// ListStarredCmd ---
type ListStarredCmd struct {
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	URLS          bool   `arg,optional,name:"urls" short:"u" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" xor:"stream" help:"Display items with their IDs"`
}

// Run (*ListStarredCmd) ---
func (s *ListStarredCmd) Run(ctx *kong.Context) error {

	if s.URLS {
		err := printStreamContentsWithURL(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}
	}

	if s.IDS {
		err := printStreamContentsWithIDs(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}
	}

	if !s.IDS && !s.URLS {
		err := printStreamContentsWithDate(s.Num, s.Order, s.ExcludeTarget, s.IncludeTarget, "user/-/state/com.google/starred")
		if err != nil {
			return err
		}
	}

	return nil
}

// WebPagesCmd ---
type WebPagesCmd struct {
	List ListWebPagesCmd `cmd:"" aliases:"ls" help:"List saved web pages"`
}

// ListWebPagesCmd ---
type ListWebPagesCmd struct {
	Num           string `arg,optional,name:"num" short:"n" help:"Specify number of items returned"`
	Order         string `arg,optional,name:"order" short:"r" help:"Specify order of items returned, i.e. newest or oldest"`
	ExcludeTarget string `arg,optional,help:"Specify stream ID of target to exclude"`
	IncludeTarget string `arg,optional,help:"Specify stream ID of target to include"`
	StreamID      string `arg,required,name:"stream-id" short:"s" help:"Specify stream ID of feed"`
	URLS          bool   `arg,optional,name:"urls" short:"u" xor:"stream" help:"Display items with their URLs"`
	IDS           bool   `arg,optional,name:"ids" short:"i" xor:"stream" help:"Display items with their IDs"`
}

// Run (*ListWebPagesCmd) ---
func (w *ListWebPagesCmd) Run(ctx *kong.Context) error {

	if w.URLS {
		err := printStreamContentsWithURL(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
			return err
		}
	} else if w.IDS {
		err := printStreamContentsWithIDs(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
			return err
		}
	} else {
		err := printStreamContentsWithDate(w.Num, w.Order, w.ExcludeTarget, w.IncludeTarget, "user/-/state/com.google/saved-web-pages")
		if err != nil {
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

	return nil
}

var cli struct {
	UserInfo UserInfoCmd `cmd:"" help:"Print Inoreader user information"`

	Login LoginCmd `cmd:"" help:"Login and initiate Oauth flow"`

	Subscription SubscriptionCmd `cmd:"" aliases:"sub" help:"Query subscriptions"`

	Tags TagsCmd `cmd:"" help:"Query tags and folders"`

	MarkItemRead MarkItemReadCmd `cmd:"" aliases:"mark" help:"Mark items as read/unread or starred/unstarred"`

	Stream StreamCmd `cmd:"" help:"Query stream contents"`

	MarkStreamRead MarkStreamReadCmd `cmd:"" help:"Mark all items in stream as read"`

	Starred StarredCmd `cmd:"" help:"Query starred articles"`

	WebPages WebPagesCmd `cmd:"" aliases:"wp" help:"Query saved web pages"`
}

func main() {

	ctx := kong.Parse(&cli,
		kong.Name("go-inoreader"),
		kong.Description("An Inoreader API client for GoLang"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			NoAppSummary: false,
			Compact:      true,
			Tree:         true,
		}))
	err := ctx.Run(&kong.Context{})
	ctx.FatalIfErrorf(err)
}
