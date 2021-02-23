# Log

### February 13, 2021
* Planned out some of the command line structure and implemented the functions for subcommands `sub all`, `sub unread`, `tag all`, `tag all -t <type>`, and `tag unread`.

### February 15, 2021
* Implemented the functions for `tag unread -t <type>`, `stream`, and `read`.
* Refactored the `sub all`, `sub unread`, `tag all`, `tag all -t <type>`, and `tag unread` subcommands. Since the stream subcommand flags modify the URL query parameters, the logic that the flags trigger will have to call `getStreamContents()` with the specified query parameters.
* Added more of the command line structure for adding and editing subscriptions and tags.
* Realized it is time to take first steps into the wonderful world of unit testing. 🎉
* Tomorrow I shall work on unit testing and implementing the command logic. 🚀

#### Resources
* [testing](https://golang.org/pkg/testing/)
* [Golang Unit Testing - Golang Docs](https://golangdocs.com/golang-unit-testing)
* [flag](https://golang.org/pkg/flag/)
* [How to use the flag package in Go](https://www.digitalocean.com/community/tutorials/how-to-use-the-flag-package-in-go)

### February 16, 2021
* Implemented the functions for `add`, `sub edit`, `tags edit`, `tags mv`, and `tags rm`.
* Combined `subscription.go`, `tags.go`, and `stream.go` into `content.go`. I will need to add comments to `content.go` to explain what each function does, and include an outline/TOC at the beginning.
* Started implementing the command logic using the flags package.
* No unit testing today; I shall implement that tomorrow after I finish the command logic.

### February 23, 2021
* Am now using the execellent [kong](https://github.com/alecthomas/kong) package to handle command line args.
* Going to hold off on unit testing until this is feature-complete 
* Made OAuth code and configuration code into one file. No reason to keep them separate.
* Added more items to [TODO.md](TODO.md).
* Implemented most of the subscription, tags, and stream functions.