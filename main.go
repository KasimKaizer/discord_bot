package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const prefix string = ">"

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	ds, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal(err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	ds.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	ds.Identify.Intents = discordgo.IntentsGuildMessages

	err = ds.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	args := strings.Split(m.Content, " ")

	if args[0] != prefix {
		return
	}

	if m.Content == "hello" {
		s.ChannelMessageSend(m.ChannelID, "world")
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}
