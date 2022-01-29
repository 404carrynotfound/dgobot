package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"strconv"
)

// Finds user current voice channel
func findChannel(session *discordgo.Session, interaction *discordgo.Interaction) string {

	for _, guild := range session.State.Guilds {
		if guild.ID == interaction.GuildID {
			for _, state := range guild.VoiceStates {
				if state.UserID == interaction.Member.User.ID {
					return state.ChannelID
				}
			}
		}
	}

	return ""
}

func validURL(value string) bool {
	_, err := url.ParseRequestURI(value)
	return err == nil
}

// Sets bot status and initialize all application commands
func ready(session *discordgo.Session, _ *discordgo.Ready) {
	// Set the playing status.
	err := session.UpdateGameStatus(0, "Serving "+strconv.Itoa(len(session.State.Guilds))+" guilds!")
	if err != nil {
		fmt.Printf("Can't set status, %s\n", err)
	}

	// Checks for unused commands and deletes them
	if commands, err := session.ApplicationCommands(session.State.User.ID, guildId); err == nil {
		found := false

		for _, l := range Commands {
			found = false

			for _, o := range commands {
				// We compare every online command with the ones locally stored, to find if a command with the same name exists
				if l.Name == o.Name {
					_, err = session.ApplicationCommandCreate(session.State.User.ID, guildId, l)
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

				_, err = session.ApplicationCommandCreate(session.State.User.ID, guildId, l)
				if err != nil {
					fmt.Printf("Cannot create '%s' command: %s\n", l.Name, err)
				}
			}
		}
	}
}
