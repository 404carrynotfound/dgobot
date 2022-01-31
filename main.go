package main

import (
	"dgobot/player"
	"fmt"
	"github.com/DisgoOrg/disgolink/dgolink"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/spotify-plugin"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

var discordToken string

const guildId = "935103396304785468"

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not find a .env file")
	}

	discordToken = os.Getenv("DISCORD_TOKEN")
}

func main() {

	//log.SetLevel(log.LevelDebug)

	if discordToken == "" {
		fmt.Println("No token provided. Please modify .env")
		return
	}

	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Printf("Error while creating Discord session %s\n", err)
		return
	}

	bot := &player.Bot{
		Link:           dgolink.New(session, lavalink.WithPlugins(spotify.New())),
		PlayerManagers: map[string]*player.Manager{},
		Session:        session,
	}

	session.AddHandler(ready)

	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.User == nil {
			if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
				handler(bot, interaction)
			}
		}
	})

	err = session.Open()
	if err != nil {
		fmt.Printf("Error while opening Discord websocket: %s\n", err)
		return
	}

	bot.RegisterNodes()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = session.Close()
	if err != nil {
		fmt.Printf("Error while closing Discord websocket: %s\n", err)
		return
	}
}
