package main

import (
	"fmt"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Plays a song from YouTube",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "song",
					Description: "YT link or song title",
					Required:    true,
				},
			},
		},
		{
			Name:        "info",
			Description: "Gives user information",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Tag user for information",
					Required:    false,
				},
			},
		},
		//{
		//	Name:        "user role",
		//	Description: "Add role to user",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionUser,
		//			Name:        "user",
		//			Description: "User",
		//			Required:    true,
		//		},
		//	},
		//},
		//{
		//	Name:        "create role",
		//	Description: "Creates role in current guild",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionString,
		//			Name:        "role name",
		//			Description: "Adds name to role",
		//			Required:    true,
		//		},
		//	},
		//},
	}
	CommandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate, bot *Bot){
		// Plays a song from spotify playlist. If it's not a valid link, it will insert into the queue the first result for the given queue
		"play": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, bot *Bot) {
			vc := findChannel(session, interaction.Interaction)

			if vc == "" {
				_, err := session.ChannelMessageSend(interaction.ChannelID, "Please provide a channel id and something to play")
				if err != nil {
					fmt.Printf("Error while sending message: %s\n", err)
				}
			}

			//link := lavalink.New(lavalink.WithUserID("935108152695881789"))

			query := interaction.Interaction.ApplicationCommandData().Options[0].StringValue()
			fmt.Println(query)
			if !validURL(query) {
				query = "ytsearch:" + query
			}

			_ = bot.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
				func(track lavalink.AudioTrack) {
					play(session, bot.Link, interaction.GuildID, vc, interaction.ChannelID, track)
				},
				func(playlist lavalink.AudioPlaylist) {
					play(session, bot.Link, interaction.GuildID, vc, interaction.ChannelID, playlist.Tracks[0])
				},
				func(tracks []lavalink.AudioTrack) {
					play(session, bot.Link, interaction.GuildID, vc, interaction.ChannelID, tracks[0])
				},
				func() {
					_, _ = session.ChannelMessageSend(interaction.ChannelID, "no matches found for: "+query)
				},
				func(ex lavalink.FriendlyException) {
					_, _ = session.ChannelMessageSend(interaction.ChannelID, "error while loading track: "+ex.Message)
				},
			))
		},

		"info": func(session *discordgo.Session, interaction *discordgo.InteractionCreate, bot *Bot) {
			sendEmbed(session, &discordgo.MessageEmbed{
				URL:         "",
				Type:        "",
				Title:       interaction.Member.User.String(),
				Description: interaction.Member.User.AvatarURL(""),
				Timestamp:   "",
				Color:       0,
				Footer:      nil,
				Image:       nil,
				Thumbnail:   nil,
				Video:       nil,
				Provider:    nil,
				Author:      nil,
				Fields:      nil,
			}, interaction.Interaction)
			fmt.Println(interaction.Member.User.ID)
		},
	}
)
