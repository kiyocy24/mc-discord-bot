package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DiscordTokenEnv     = "DISCORD_TOKEN"
	DiscordClientIDEnv  = "DISCORD_CLIENT_ID"
	DiscordChannelIDEnv = "DISCORD_CHANNEL_ID"

	GCPAPIkeyEnv       = "GOOGLE_APPLICATION_CREDENTIALS"
	GCPProjectIDEnv    = "GCP_PROJECT_ID"
	GCPStartSubIDEnv   = "GCP_START_SUB_ID"
	GCPStopSubIDEnv    = "GCP_STOP_SUB_ID"
	GCPRestartSubIDEnv = "GCP_RESTART_SUB_ID"

	Command       = "/"
	StartCommmand = Command + "start"
	StopCommand   = Command + "stop"
	StatusCommand = Command + "status"
)

func main() {
	fmt.Println("Bot start")

	// init env
	token := os.Getenv(DiscordTokenEnv)
	if token == "" {
		fmt.Println(envError(DiscordTokenEnv))
	}
	fmt.Println("token : ", token)

	client := os.Getenv(DiscordClientIDEnv)
	if client == "" {
		fmt.Println(envError(DiscordClientIDEnv))
	}
	fmt.Println("client : ", client)

	channelID := os.Getenv(DiscordChannelIDEnv)
	if channelID == "" {
		fmt.Println(envError(DiscordChannelIDEnv))
	}
	fmt.Println("channel id : ", channelID)

	// new discord token
	dg, err := discordgo.New("Bot " + token)
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
		sendMessage, err := s.ChannelMessageSend(m.ChannelID, "Server starting...")
		if err != nil {
			fmt.Printf("discord error : %v", err)
		}

		// start function
		subID := os.Getenv(GCPStartSubIDEnv)
		if subID == "" {
			fmt.Println(envError(GCPStartSubIDEnv))
			return
		}

		err = pullMsgsSync(bytes.NewBufferString("Start instance"), subID)
		if err != nil {
			s.ChannelMessageEdit(m.ChannelID, sendMessage.ID, "Failed: Server did not started up.")
			fmt.Printf("pullMesgsSync error : %v\n", err)
			return
		}

		s.ChannelMessageSend(m.ChannelID, "Success! Server started up.")
		fmt.Printf("Server started up.")
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

func pullMsgsSync(w io.Writer, subID string) error {
	projectID := os.Getenv(GCPProjectIDEnv)
	if projectID == "" {
		fmt.Println(envError(GCPProjectIDEnv))
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)
	exists, err := sub.Exists(ctx)
	if err != nil {
		fmt.Println("test")
		return fmt.Errorf("%v", err)
	}
	if !exists {
		return fmt.Errorf("subscription is not exist. sub id: %v", subID)
	}

	sub.ReceiveSettings.Synchronous = true

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	cm := make(chan *pubsub.Message)
	defer close(cm)

	go func() {
		for msg := range cm {
			fmt.Fprintf(w, "Got message :%q\n", string(msg.Data))
			msg.Ack()
		}
	}()

	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		cm <- msg
	})
	if err != nil && status.Code(err) != codes.Canceled {
		return fmt.Errorf("Receive: %v", err)
	}

	return nil
}

func envError(s string) error {
	return errors.New("Error: " + s + " env is not set")
}
