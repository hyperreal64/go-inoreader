# go-inoreader
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperreal64/go-inoreader)](https://goreportcard.com/report/github.com/hyperreal64/go-inoreader)

WORK IN PROGRESS 🚧: An unofficial Inoreader API client

We need to create a new application on Inoreader under Preferences > Developer. Set the redirect URI to `http://localhost:8081/oauth/redirect` and scope to `Read and write`. We will then get an App ID and App Key. Save the App ID and App Key aa JSON items to the configuration file. On Unix/Linux: `~/.local/share/go-inoreader.json`. On Windows: `$env:APPDATA\go-inoreader.json`.

```json
{
  "app_id": <your app id>,
  "app_key": <your app key>
}
```

Install using Go's builtin package manager:
```bash
go get -v github.com/hyperreal64/go-inoreader
```

## Usage

We need to authorize with the Inoreader API first:
```go
package main

import "github.com/hyperreal64/go-inoreader/config"

func main() {
  	config.Init()
}
```

Now run `go run main.go` to initiate the OAuth flow.

### Example: Subscription list

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hyperreal64/go-inoreader/config"
	"github.com/hyperreal64/go-inoreader/subscription"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	rc := config.Oauth2RestyClient(ctx)
	defer cancel()
	
	sublist, err := subscription.GetSubscriptionList(rc)
	if err != nil {
		log.Fatalln(err)
	}
	
	for _, v := range sublist.Subscriptions {
		fmt.Printf("%s (%s)\n", v.Title, v.URL)
	}
}
```

Output would be something like this:
```
OSNews (http://www.osnews.com/files/recent.xml)
Planet Python (http://planet.python.org/rss10.xml)
Opensource.com (https://opensource.com/feed)
The GoLang Blog (https://blog.golang.org/feed.atom)
Going Go Programming (https://www.goinggo.net/index.xml)
Go Time (https://changelog.com/gotime/feed)
The Changelog (https://changelog.com/podcast/feed)
Command Line Heroes (https://feeds.pacific-content.com/commandlineheroes)
freeCodeCamp.org News (https://www.freecodecamp.org/news/rss/)
Enable Sysadmin (https://www.redhat.com/sysadmin/rss.xml)
Blogs on Drew DeVault's blog (https://drewdevault.com/blog/index.xml)
```
