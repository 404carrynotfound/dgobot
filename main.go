package main

import (
	"fmt"
	"github.com/DisgoOrg/disgolink/dgolink"
	"github.com/DisgoOrg/log"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	discordToken string
	//spotifySecret string
	//spotifyID     string
	prefix string
	//spotifyClient spotify.Client
)

type Bot struct {
	Link *dgolink.Link
}

const guildId = "935103396304785468"

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not find a .env file")
	}

	discordToken = os.Getenv("DISCORD_TOKEN")
	//spotifyID = os.Getenv("SPOTIFY_ID")
	//spotifySecret = os.Getenv("SPOTIFY_SECRET")

	//// Spotify credentials
	//spotifyConfig := &clientcredentials.Config{
	//	ClientID:     spotifyID,
	//	ClientSecret: spotifySecret,
	//	TokenURL:     spotify.TokenURL,
	//}
	//
	//// Check spotify token and create spotify client
	//spotifyToken, err := spotifyConfig.Token(context.Background())
	//if err != nil {
	//	fmt.Printf("Spotify: couldn't get token: %s\n", err)
	//}
	//
	//spotifyClient = spotify.Authenticator{}.NewClient(spotifyToken)
}

func main() {

	log.SetLevel(log.LevelDebug)

	if discordToken == "" {
		fmt.Println("No token provided. Please modify .env")
		return
	}

	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Printf("Error while creating Discord session %s\n", err)
		return
	}

	bot := &Bot{
		Link: dgolink.New(session),
	}

	session.AddHandler(ready)

	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.User == nil {
			if handler, ok := CommandHandlers[interaction.ApplicationCommandData().Name]; ok {
				handler(session, interaction, bot)
			}
		}
	})

	err = session.Open()
	if err != nil {
		fmt.Printf("Error while opening Discord websocket: %s\n", err)
		return
	}

	bot.registerNodes()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = session.Close()
	if err != nil {
		fmt.Printf("Error while closing Discord websocket: %s\n", err)
		return
	}
}

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	// Set the playing status.
	err := s.UpdateGameStatus(0, "Serving "+strconv.Itoa(len(s.State.Guilds))+" guilds!")
	if err != nil {
		fmt.Printf("Can't set status, %s\n", err)
	}

	// Checks for unused commands and deletes them
	if cmds, err := s.ApplicationCommands(s.State.User.ID, guildId); err == nil {
		found := false

		for _, l := range Commands {
			found = false

			for _, o := range cmds {
				// We compare every online command with the ones locally stored, to find if a command with the same name exists
				if l.Name == o.Name {
					_, err = s.ApplicationCommandCreate(s.State.User.ID, guildId, l)
					if err != nil {
						fmt.Printf("Cannot create '%s' command: %s\n", l.Name, err)
					}

					found = true
					break
				}

			}
			// If we didn't found a match for the locally stored command, it means the command is new. We register it
			if !found {
				fmt.Printf("Registering new command %s\n", l.Name)

				_, err = s.ApplicationCommandCreate(s.State.User.ID, guildId, l)
				if err != nil {
					fmt.Printf("Cannot create '%s' command: %s", l.Name, err)
				}
			}
		}
	}
}
