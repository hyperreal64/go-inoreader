# go-inoreader

WORK IN PROGRESS 🚧: An unofficial Inoreader API client for Go

The general guidelines I use for coding this are:

🚴 Stop [bike-shedding](https://en.wikipedia.org/wiki/Law_of_triviality)

👍 Solve the real problem

💩 First, write shxtty code

🌟 Figure out how to make it better

--Inanc Gumus [@inancgumus](https://twitter.com/inancgumus)

To use this, you need to create a new application on Inoreader under Preferences > Developer. You will be give an App ID and App Secret for your new app. On Unix/Linux, set these values in the following environment variables:
```bash
export INO_APP_ID=<app id>
export INO_APP_SEC=<app secret>
```

You may also save them to the shell's profile so that they are sourced on startup.

For Windows, environment variables can be set by going to Explorer.exe > This PC > Right click and select Properties > Advanced system settings > Environment Variables. Create a new User environment variable for INO_APP_ID and INO_APP_SEC and set their values to App ID and App Secret.