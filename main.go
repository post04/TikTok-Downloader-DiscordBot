package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tiktok "./tiktok"
	"github.com/bwmarrin/discordgo"
)

var (
	prefix, token string
	embed         = &discordgo.MessageEmbed{}
)

type config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

func loadconfig() (string, string) {
	return Token, Prefix
}

func messageCreate(session *discordgo.Session, msg *discordgo.MessageCreate) {
	channel, err := session.Channel(msg.ChannelID)
	if err != nil {
		return
	}
	if channel.Type != 1 {
		return
	}
	var args = strings.Split(msg.Content, " ")
	var command string
	if len(args) < 1 {
		return
	}
	if !strings.HasPrefix(msg.Content, prefix) {
		return
	}
	command = strings.Replace(args[0], prefix, "", 1)
	args = args[1:]

	command = strings.ToLower(command)
	if command == "dltiktok" {
		if len(args) < 1 {
			session.ChannelMessageSend(msg.ChannelID, "Please provide valid input, Like this: \nhttps://vm.tiktok.com/ZMJvnDFWr/")
			return
		}
		fmt.Println("User", msg.Author.Username+"#"+msg.Author.Discriminator, "used dltiktok command!")
		if !strings.Contains(args[0], ".") {
			session.ChannelMessageSend(msg.ChannelID, "Please provide valid input, Like this: \nhttps://vm.tiktok.com/ZMJvnDFWr/")
			return
		}
		if !strings.HasPrefix(args[0], "https://") || strings.Split(args[0], ".")[1] != "tiktok" {
			session.ChannelMessageSend(msg.ChannelID, "Please provide valid input, Like this: \nhttps://vm.tiktok.com/ZMJvnDFWr/")
			return
		}
		var tiktokvideo = args[0]
		if !strings.HasSuffix(args[0], "/") {
			tiktokvideo += "/"
		}
		video := tiktok.Video{URL: tiktokvideo, BaseDIR: "./tiktok", Proxy: ""}
		err := video.FetchInfo()
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Error downloading tiktok!\n```"+fmt.Sprint(err)+"```")
			return
		}
		filedata, err := video.Download()
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Error downloading tiktok!\n```"+fmt.Sprint(err)+"```")
			return
		}

		_, err = session.ChannelFileSend(msg.ChannelID, "test.mp4", filedata)
		if err != nil {
			session.ChannelMessageSend(msg.ChannelID, "Error downloading tiktok!\n```"+fmt.Sprint(err)+"```")
			return
		}
		return
	}
	if command == "help" {
		fmt.Println("User", msg.Author.Username+"#"+msg.Author.Discriminator, "used help command!")
		session.ChannelMessageSendEmbed(msg.ChannelID, embed)
		return
	}

}

func ready(session *discordgo.Session, event *discordgo.Ready) {
	session.UpdateStatus(0, "Please use "+prefix+"dltiktok in my dms or "+prefix+"help!")
	fmt.Println("Logged in as: " + event.User.Username + "#" + event.User.Discriminator)
}

func main() {
	token, prefix = loadconfig()
	embed = &discordgo.MessageEmbed{
		Title:       "Tiktok video downloader!",
		Description: "Please use the command `" + prefix + "dltiktok` in my dms to download a tiktok of your choice! Make sure the tiktok is using the link that looks like this: `https://vm.tiktok.com/ZMJvnDFWr/` \nor \n`https://www.tiktok.com/@msbutterworththicknrich/video/6917677526236089606`\nExample: `" + prefix + "dltiktok https://vm.tiktok.com/ZMJvnDFWr/`\n\nInvite: [Click me!](https://discord.com/api/oauth2/authorize?client_id=798439375712813126&permissions=116736&scope=bot)\nSupport Server: [Click me!](https://discord.gg/vBxmKKE)\nGithub link: [Click me!](https://github.com/postrequest69/TikTok-Downloader-DiscordBot)",
	}
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session")
		return
	}
	bot.AddHandler(ready)
	bot.AddHandler(messageCreate)
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,")
		return
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	_ = bot.Close()
}
