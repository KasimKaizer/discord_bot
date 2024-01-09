package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	prefix string
	token  string
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error while loading the env file", err)
	}

	log.Println("env loaded successfully")
	prefix = os.Getenv("PREFIX")
	token = os.Getenv("BOT_TOKEN")
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

	message := m.Content[1:]

	if message == "cat" || message == "dog" || message == "meme" {
		var url string

		if message == "cat" {
			url = getCat()
		}

		if message == "dog" {
			url = getDog()
		}

		if message == "meme" {
			url = getMeme()
		}

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
		s.ChannelMessageSend(m.ChannelID, url)

	}

	if message == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func getCat() string {
	resp, err := http.Get("https://api.thecatapi.com/v1/images/search?size=full")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp)
	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var result []map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result[0]["url"].(string))
	return result[0]["url"].(string)
}

func getDog() string {
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		log.Fatal(err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result["message"].(string))
	return result["message"].(string)
}

func getMeme() string {
	resp, err := http.Get("https://meme-api.com/gimme")
	if err != nil {
		log.Fatal(err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result map[string]interface{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result["url"].(string))
	return result["url"].(string)
}
