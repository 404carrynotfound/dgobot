package main

import (
	"dgobot/interaction"
	"dgobot/player"
	"fmt"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
	"strconv"
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
		{
			Name:        "kick",
			Description: "Kick user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Tag user to be kicked",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for this ban",
					Required:    true,
				},
			},
		},
		{
			Name:        "ban",
			Description: "Ban user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "Tag user to be banned",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "The reason for this ban",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "days",
					Description: "The number of days of previous comments to delete.",
					Required:    true,
				},
			},
		},
		{
			Name:        "unban",
			Description: "Unban user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "user",
					Description: "Unban user",
					Required:    true,
				},
			},
		},
		//{
		//	Name:        "mute",
		//	Description: "mute user",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionUser,
		//			Name:        "user",
		//			Description: "mute user",
		//			Required:    true,
		//		},
		//	},
		//},
		//{
		//	Name:        "deaf",
		//	Description: "deaf user",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:        discordgo.ApplicationCommandOptionUser,
		//			Name:        "user",
		//			Description: "deaf user",
		//			Required:    true,
		//		},
		//	},
		//},
		//{
		//	Name:        "purge",
		//	Description: "delete last n messages",
		//	Options: []*discordgo.ApplicationCommandOption{
		//		{
		//			Type:     discordgo.ApplicationCommandOptionInteger,
		//			Name:     "count",
		//			Required: true,
		//		},
		//	},
		//},
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
					bot.Play(vc, interaction.Interaction, "", track)
				},
				func(playlist lavalink.AudioPlaylist) {
					bot.Play(vc, interaction.Interaction, playlist.Info.Name, playlist.Tracks...)
				},
				func(tracks []lavalink.AudioTrack) {
					bot.Play(vc, interaction.Interaction, "", tracks[0])
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
			bot.Skip(interaction.Interaction)
		},

		"stop": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Stop(interaction.Interaction)
			err := bot.Session.ChannelVoiceJoinManual(guildId, "", false, false)
			if err != nil {
				fmt.Printf("Error when leaving channel: %s\n", err)
				return
			}
		},

		"current": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Current(interaction.Interaction)
		},

		"pause": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Pause(interaction.Interaction)
		},

		"resume": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			bot.Resume(interaction.Interaction)
		},

		"whois": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			var user *discordgo.User
			if len(interaction.Interaction.ApplicationCommandData().Options) != 0 {
				user = interaction.Interaction.ApplicationCommandData().Options[0].UserValue(bot.Session)
			}

			if user == nil {
				user = interaction.Member.User
			}

			member, _ := bot.Session.GuildMember(interaction.GuildID, user.ID)

			interactions.SendEmbedInteraction(bot.Session, &discordgo.MessageEmbed{
				Title:     member.User.String(),
				Thumbnail: &discordgo.MessageEmbedThumbnail{URL: member.User.AvatarURL("")},
				Author: &discordgo.MessageEmbedAuthor{
					Name:    member.User.String(),
					IconURL: member.User.AvatarURL(""),
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Joined",
						Value:  member.JoinedAt.Local().Format("Mon Jan 2, 2006 15:04"),
						Inline: false,
					},
					{
						Name:   "Muted",
						Value:  strconv.FormatBool(member.Mute),
						Inline: true,
					},
					{
						Name:   "Deaf",
						Value:  strconv.FormatBool(member.Deaf),
						Inline: true,
					},
				},
			}, interaction.Interaction)
		},

		"kick": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			reason := interaction.ApplicationCommandData().Options[1].StringValue()

			err := bot.Session.GuildMemberDeleteWithReason(interaction.GuildID, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).ID, reason)
			if err != nil {
				interactions.SendAndDeleteInteraction(bot.Session, "Error has occurred", interaction.Interaction, time.Second*5)
				fmt.Printf("Error with kicking member from guild: %s\n", err)
				return
			}
			interactions.SendMessageInteraction(bot.Session, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).Mention()+" is kicked", interaction.Interaction)
		},

		"ban": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			reason := interaction.ApplicationCommandData().Options[1].StringValue()
			days := interaction.ApplicationCommandData().Options[2].IntValue()

			err := bot.Session.GuildBanCreateWithReason(interaction.GuildID, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).ID, reason, int(days))
			if err != nil {
				interactions.SendAndDeleteInteraction(bot.Session, "Error has occurred", interaction.Interaction, time.Second*5)
				fmt.Printf("Error with banning member from guild: %s\n", err)
				return
			}
			interactions.SendMessageInteraction(bot.Session, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).Mention()+" is banned", interaction.Interaction)
		},

		"unban": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
			err := bot.Session.GuildBanDelete(interaction.GuildID, interaction.ApplicationCommandData().Options[0].StringValue())
			if err != nil {
				interactions.SendAndDeleteInteraction(bot.Session, "Invalid id", interaction.Interaction, time.Second*5)
				fmt.Printf("Error with banning member from guild: %s\n", err)
				return
			}
			interactions.SendMessageInteraction(bot.Session, "<@"+interaction.ApplicationCommandData().Options[0].StringValue()+"> is unbanned", interaction.Interaction)
		},
	}
)
