package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	Command       = "/ "
	StartCommmand = Command + "start"
	StopCommand   = Command + "stop"
	StatusCommand = Command + "status"
)

var (
	TokenIdError   = errors.New("Token ID is invalid")
	ClientIdError  = errors.New("Client ID is invalid")
	ChannelIdError = errors.New("Channel ID is invalid")
)

var (
	Token     string
	ClientId  string
	ChannelId string
)

func initValue() error {
	Token = os.Getenv("TOKEN")
	if Token == "" {
		return TokenIdError
	}

	ClientId = os.Getenv("CLIENT_ID")
	if ClientId == "" {
		return ClientIdError
	}

	ChannelId = os.Getenv("CHANNEL_ID")
	if Token == "" {
		return ChannelIdError
	}

	return nil

}

func main() {
	fmt.Println("Bot start")
	err := initValue()
	if err != nil {
		fmt.Println("Init value error.", err)
		return
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Login error.", err)
		return
	}

	dg.AddHandler(messageCreate)
	err = dg.Open()

	if err != nil {
		fmt.Println("error opening connection,", err)
	}
	defer dg.Close()

	fmt.Println("bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == StartCommmand {
		s.ChannelMessageSend(m.ChannelID, "Server starting...")
		// start function
	}

	if m.Content == StopCommand {
		s.ChannelMessageSend(m.ChannelID, "Server stopping...")
		// stop function
	}

	if m.Content == StatusCommand {
		s.ChannelMessageSend(m.ChannelID, "Server status is ...")
		// status function
	}
}
