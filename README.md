# go-inoreader
[![Go Report Card](https://goreportcard.com/badge/github.com/hyperreal64/go-inoreader)](https://goreportcard.com/report/github.com/hyperreal64/go-inoreader)

WORK IN PROGRESS 🚧: An unofficial Inoreader API client for Go

The general guidelines I use for coding this are:

🚴 Stop [bike-shedding](https://en.wikipedia.org/wiki/Law_of_triviality)

👍 Solve the real problem

💩 First, write shxtty code

🌟 Figure out how to make it better

--Inanc Gumus [@inancgumus](https://twitter.com/inancgumus)

To use this, you need to create a new application on Inoreader under Preferences > Developer. Set the redirect URI to `http://localhost:8081/oauth/redirect` and scope to `Read and write`. You will get an App ID and App Key which you need to save to the configuration file. On Unix/Linux, open `~/.local/share/go-inoreader.json`. On Windows, open `$HOME\AppData\Local\go-inoreader.json`. Add your App ID and App Key as JSON items to this file.
```json
{
  "app_id": <your app id>,
  "app_key": <your app key>
}
```

Save and close the file. Clone this repo and build the binary:
```bash
git clone https://github.com/hyperreal64/go-inoreader
cd go-inoreader
go build .
```

To see a list of commands:
```bash
./go-inoreader --help
```

To initiate the OAuth flow:
```bash
./go-inoreader login
```

Point your browser at [http://localhost:8081](http://localhost:8081) and click 'Authorize'.

Requires Go 1.16+
