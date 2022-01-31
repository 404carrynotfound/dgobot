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

// Finds if user have a specific role
func findUserRole(member *discordgo.Member, role *discordgo.Role) bool {
	for _, memberRole := range member.Roles {
		if memberRole == role.ID {
			return true
		}
	}
	return false
}

// Validate URL
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
	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, guildId, Commands)
	if err != nil {
		fmt.Printf("Error while loading slash commands: %s\n", err)
	}
}
