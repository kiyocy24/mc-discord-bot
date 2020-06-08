package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token     string
	ClientId  string
	ChannelId string
)

func initValue() error {
	Token = os.Getenv("TOKEN")
	if Token == "" {
		return errors.New("Token ID is empty")
	}

	ClientId = os.Getenv("CLIENT_ID")
	if ClientId == "" {
		return errors.New("Client ID is empty")
	}

	ChannelId = os.Getenv("CHANNEL_ID")
	if Token == "" {
		return errors.New("Channel ID is empty")
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

	if m.Content == "/mc start" {
		s.ChannelMessageSend(m.ChannelID, "Server starting...")
	}

	if m.Content == "/mc stop" {
		s.ChannelMessageSend(m.ChannelID, "Server stopping...")
	}

	if m.Content == "/mc status" {
		s.ChannelMessageSend(m.ChannelID, "Server status is ...")
	}
}
