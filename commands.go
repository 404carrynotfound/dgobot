package main

import (
	"dgobot/interaction"
	"dgobot/player"
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
	CommandHandlers = map[string]func(bot *player.Bot, interaction *discordgo.InteractionCreate){
		// Plays a song from spotify playlist. If it's not a valid link, it will insert into the queue the first result for the given queue
		"play": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			vc := findChannel(bot.Session, interaction.Interaction)
			if vc == "" {
				interactions.SendAndDeleteInteraction(bot.Session, "Please join a voice channel.", interaction.Interaction, time.Second*5)
				return
			}

			query := interaction.Interaction.ApplicationCommandData().Options[0].StringValue()
			if !validURL(query) {
				query = "ytsearch:" + query
			}

			err := bot.Link.BestRestClient().LoadItemHandler(query, lavalink.NewResultHandler(
				func(track lavalink.AudioTrack) {
					bot.Play(bot.Session, vc, interaction.Interaction, "", track)
				},
				func(playlist lavalink.AudioPlaylist) {
					bot.Play(bot.Session, vc, interaction.Interaction, playlist.Info.Name, playlist.Tracks...)
				},
				func(tracks []lavalink.AudioTrack) {
					bot.Play(bot.Session, vc, interaction.Interaction, "", tracks[0])
				},
				func() {
					_, _ = bot.Session.ChannelMessageSend(interaction.ChannelID, "no matches found for: "+query)
				},
				func(ex lavalink.FriendlyException) {
					_, _ = bot.Session.ChannelMessageSend(interaction.ChannelID, "error while loading track: "+ex.Message)
				},
			))
			if err != nil {
				fmt.Printf("Error while player loading: %s\n", err)
				return
			}
		},

		"skip": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Skip(bot.Session, interaction.Interaction)
		},

		"stop": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Stop(bot.Session, interaction.Interaction)
			err := bot.Session.ChannelVoiceJoinManual(guildId, "", false, false)
			if err != nil {
				fmt.Printf("Error when leaving channel: %s\n", err)
				return
			}
		},

		"current": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Current(bot.Session, interaction.Interaction)
		},

		"pause": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Pause(bot.Session, interaction.Interaction)
		},

		"resume": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Resume(bot.Session, interaction.Interaction)
		},

		"whois": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			var user *discordgo.User
			if len(interaction.Interaction.ApplicationCommandData().Options) != 0 {
				user = interaction.Interaction.ApplicationCommandData().Options[0].UserValue(bot.Session)
			}

			if user == nil {
				user = interaction.Member.User
			}

			interactions.SendEmbedInteraction(bot.Session, &discordgo.MessageEmbed{
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
