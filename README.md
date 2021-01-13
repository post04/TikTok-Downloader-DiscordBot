# TikTok-Downloader-DiscordBot
A discord bot that allows you to download TikTok videos with or without downloads enabled, it's dm only right now.

To fillout the config you need to edit config.go, I did this because when compiling to a linux binary and running with a service, it seemed to not be able to pickup the config.json file.

# To run

Make sure you have [go](https://golang.org/doc/install) installed!
```
$ git clone https://github.com/postrequest69/TikTok-Downloader-DiscordBot.git
$ cd TikTok-Downloader-DiscordBot
$ go build 
$ run the exe...
```

# TODO
Accept any tiktok link.
Make it not dm only but with a toggleable guild config for dm only.
