package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	prefix   string
	token    string
	nightApi string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while loading the env file", err)
	}

	log.Println("env loaded successfully")
	prefix = os.Getenv("PREFIX")
	token = os.Getenv("BOT_TOKEN")
	nightApi = os.Getenv("NIGHT_API")

}

func main() {
	ds, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	ds.AddHandler(messageCreate)

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

	if m.Content[:1] != prefix {
		return
	}
	message := strings.Split(m.Content[1:], " ")

	if message[0] == "pet" {

		url := getPet(strings.ToLower(message[1]))

		// TODO:
		// In future if you want to embed the image, then the code for it is here.

		// imgResp, err := http.Get(url)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// defer imgResp.Body.Close()

		// attachment := discordgo.File{
		// 	Name:   "cute_pet.jpg",
		// 	Reader: imgResp.Body,
		// }
		// s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		// 	Files: []*discordgo.File{&attachment},
		// })

		if url == "" {
			s.ChannelMessageSend(m.ChannelID, "Invalid command / Some error occurred")
		} else {
			s.ChannelMessageSend(m.ChannelID, url)
		}
	}

	if message[0] == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func getPet(pet string) string {

	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.night-api.com/images/animals/%s", pet), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("authorization", nightApi)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonResponse map[string]interface{}

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		log.Fatal(err)
	}

	// We are taking the easy way out here.
	// TODO: to something better here to detect and handle error
	if jsonResponse["status"].(float64) != 200 {
		return ""
	}

	content := jsonResponse["content"].(map[string]interface{})

	return content["url"].(string)
}
