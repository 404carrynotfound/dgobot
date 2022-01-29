package main

import (
	"dgobot/interaction"
	"fmt"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
	"time"
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
			Name:        "skip",
			Description: "skips current song",
		},
		{
			Name:        "stop",
			Description: "stop player",
		},
		{
			Name:        "current",
			Description: "current song",
		},
		{
			Name:        "pause",
			Description: "pauses player",
		},
		{
			Name:        "resume",
			Description: "resumes player",
		},
		{
			Name:        "whois",
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
	CommandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
		// Plays a song from spotify playlist. If it's not a valid link, it will insert into the queue the first result for the given queue
		"play": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			vc := findChannel(session, interaction.Interaction)
			if vc == "" {
				interactions.SendAndDeleteInteraction(session, "Please join a voice channel.", interaction.Interaction, time.Second*5)
				return
			}

			query := interaction.Interaction.ApplicationCommandData().Options[0].StringValue()
			if !validURL(query) {
				query = "ytsearch:" + query
			}

			err := botLink.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
				func(track lavalink.AudioTrack) {
					botLink.Play(session, vc, interaction.Interaction, "", track)
				},
				func(playlist lavalink.AudioPlaylist) {
					botLink.Play(session, vc, interaction.Interaction, playlist.Info.Name, playlist.Tracks...)
				},
				func(tracks []lavalink.AudioTrack) {
					botLink.Play(session, vc, interaction.Interaction, "", tracks[0])
				},
				func() {
					_, _ = session.ChannelMessageSend(interaction.ChannelID, "no matches found for: "+query)
				},
				func(ex lavalink.FriendlyException) {
					_, _ = session.ChannelMessageSend(interaction.ChannelID, "error while loading track: "+ex.Message)
				},
			))
			if err != nil {
				fmt.Printf("Error while player loading: %s\n", err)
				return
			}
		},

		"skip": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			botLink.Skip(session, interaction.Interaction)
		},

		"stop": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			botLink.Stop(session, interaction.Interaction)
			err := session.ChannelVoiceJoinManual(guildId, "", false, false)
			if err != nil {
				fmt.Printf("Error when leaving channel: %s\n", err)
				return
			}
		},

		"current": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			botLink.Current(session, interaction.Interaction)
		},

		"pause": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			botLink.Pause(session, interaction.Interaction)
		},

		"resume": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			botLink.Resume(session, interaction.Interaction)
		},

		"whois": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
			var user *discordgo.User
			if len(interaction.Interaction.ApplicationCommandData().Options) != 0 {
				user = interaction.Interaction.ApplicationCommandData().Options[0].UserValue(session)
			}

			if user == nil {
				user = interaction.Member.User
			}

			interactions.SendEmbedInteraction(session, &discordgo.MessageEmbed{
				URL:         "",
				Type:        "",
				Title:       user.String(),
				Description: "",
				Timestamp:   "",
				Color:       0,
				Footer:      nil,
				Image:       nil,
				Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: user.AvatarURL("")},
				Video:       nil,
				Provider:    nil,
				Author:      nil,
				Fields:      nil,
			}, interaction.Interaction)
		},
	}
)
