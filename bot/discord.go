package bot

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/stevegood/flare/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Discord struct {
	token string
}

func (d *Discord) Connect() error {
	session, err := discordgo.New("Bot " + d.token)
	if err != nil {
		log.Println("Error when creating new Discord session", err)
		return err
	}

	session.AddHandler(newMessageHandler)

	err = session.Open()
	if err != nil {
		log.Println("Error when opening connection to Discord", err)
		return err
	}

	_ = session.UpdateStatus(0, "Solar San Diego helper bot")

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Flare is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	_ = session.Close()

	return nil
}

func newMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!gonk" {
		gonkContent := Gonk()
		messageSendEmbed(s, m, &gonkContent)
	}

	if m.Content == "!lumpy" {
		randomLumpyURL := Lumpy()

		messageSendEmbed(s, m, &discordgo.MessageEmbed{
			Image: &discordgo.MessageEmbedImage{
				URL: randomLumpyURL,
			},
		})
	}
}

func messageSendEmbed(s *discordgo.Session, m *discordgo.MessageCreate, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Embedded item as follows")
		fmt.Printf("%d Fields\n", len(embed.Fields))

		b, err := json.Marshal(embed)
		if err != nil {
			fmt.Println("Could not unmarshal data")
			fmt.Println(err)
		}

		fmt.Println(string(b))

		outputError := Error("Failed to render", utils.WithTemplate("There was a problem rendering the %s \"%s\"", embed.Author.Name, embed.Title))
		_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &outputError)
	}
}

func NewDiscord(token string) Discord {
	discordBot := Discord{token: token}
	return discordBot
}

// Error returns an discordgo.MessageEmbed with a level of Error
func Error(title, description string) discordgo.MessageEmbed {
	return newEmbed(0xE84A4A, title, description)
}

// Info returns an discordgo.MessageEmbed with a level of Info
func Info(title, description string) discordgo.MessageEmbed {
	return newEmbed(0xF2E82B, title, description)
}

func newEmbed(color int, title, description string) discordgo.MessageEmbed {
	return discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
	}
}
