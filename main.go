package main

import (
	"discord-bot/gcp"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	DiscordTokenEnv     = "DISCORD_TOKEN"
	DiscordClientIDEnv  = "DISCORD_CLIENT_ID"
	DiscordChannelIDEnv = "DISCORD_CHANNEL_ID"

	GCPProjectIDEnv    = "GCP_PROJECT_ID"
	GCPStartTopicIDEnv = "GCP_START_TOPIC_ID"
	GCPStopTopicIDEnv  = "GCP_STOP_TOPIC_ID"

	Command       = "/"
	StartCommand  = Command + "start"
	StopCommand   = Command + "stop"
	StatusCommand = Command + "status"
)

var logInfo *log.Logger
var logError *log.Logger
var logFatal *log.Logger

func init() {
	logInfo = log.New(os.Stdout, "[INFO]", log.LstdFlags)
	logError = log.New(os.Stdout, "[ERROR]", log.LstdFlags)
	logFatal = log.New(os.Stdout, "[FATAL]", log.LstdFlags)
}

func main() {
	logInfo.Println("Bot start.")
	defer logInfo.Println("Bot finish.")

	// init env
	token := os.Getenv(DiscordTokenEnv)
	logInfo.Println("token : ", token)

	client := os.Getenv(DiscordClientIDEnv)
	logInfo.Println("client : ", client)

	channelID := os.Getenv(DiscordChannelIDEnv)
	logInfo.Println("channel id : ", channelID)

	// new discord token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logError.Println("Login error.", err)
		return
	}

	dg.AddHandler(messageCreate)
	err = dg.Open()

	if err != nil {
		logError.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	logInfo.Println("bot is now running. Press CTRL-C to exit.")
	logInfo.Println("-------------------------------------------")
	sc := make(chan os.Signal)
	signal.Notify(sc, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, os.Kill)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == StartCommand || m.Content == StopCommand {
		logInfo.Printf("id: %v, name: %v, command: %v", m.Author.ID, m.Author.Username, m.Content)
	}

	projectID := os.Getenv(GCPProjectIDEnv)

	switch m.Content {
	case StartCommand:
		sendMessage, err := s.ChannelMessageSend(m.ChannelID, "Server starting...")
		if err != nil {
			logError.Printf("discord error : %v", err)
		}
		topicID := os.Getenv(GCPStartTopicIDEnv)

		serverID, err := gcp.Publish(projectID, topicID, "start")
		logInfo.Printf("Published a message; msg ID: %v\n", serverID)

		if err != nil {
			s.ChannelMessage(m.ChannelID, "Failed: Server did not started up.\n"+err.Error())
			logError.Printf("pubsub publish error : %v\n", err)
			return
		}

		time.Sleep(time.Second * 10)
		s.ChannelMessageEdit(m.ChannelID, sendMessage.ID, "Success! Server started up.")
		logInfo.Printf("Server started up.")
		break

	case StopCommand:
		sendMessage, err := s.ChannelMessageSend(m.ChannelID, "Server stopping...")
		if err != nil {
			logError.Printf("discord error : %v", err)
		}
		topicID := os.Getenv(GCPStopTopicIDEnv)

		serverID, err := gcp.Publish(projectID, topicID, "start")
		logInfo.Printf("Published a message; msg ID: %v\n", serverID)

		if err != nil {
			s.ChannelMessage(m.ChannelID, "Failed: Server did not stop.\n"+err.Error())
			logError.Printf("pubsub publish error : %v\n", err)
			return
		}

		time.Sleep(time.Second * 10)
		s.ChannelMessageEdit(m.ChannelID, sendMessage.ID, "Success! Server stopped.")
		logInfo.Printf("Server stopped.")
		break

	}
}
