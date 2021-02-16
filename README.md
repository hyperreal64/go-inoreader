# go-inoreader

WORK IN PROGRESS: An unofficial Inoreader API client for Go

The general guidelines I use for coding this are:

🚴 Stop [bike-shedding](https://en.wikipedia.org/wiki/Law_of_triviality)

👍 Solve the real problem

💩 First, write shxtty code

🌟 Figure out how to make it better

--Inanc Gumus [@inancgumus](https://twitter.com/inancgumus)

To use this, you need to create a new application on Inoreader under Preferences > Developer. On your system, set environment variables for INOREADER_CLIENT_ID and INOREADER_CLIENT_SECRET using the App ID and App Secret, respectively. On Unix/Linux shells, this can be done with the following commands:
```bash
export INOREADER_CLIENT_ID=<app id>
export INOREADER_CLIENT_SECRET=<app secret>
```

You may also save them to the shell's profile so that they are sourced on startup.

For Windows, environment variables can be set by going to Explorer.exe > This PC > Right click and select Properties > Advanced system settings > Environment Variables. Create a new User environment variable for INOREADER_CLIENT_ID and INOREADER_CLIENT_SECRET and set their values to App ID and App Secret.