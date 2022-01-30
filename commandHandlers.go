package main

import (
	interactions "dgobot/interaction"
	"dgobot/player"
	"fmt"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"time"
)

var CommandHandlers = map[string]func(bot *player.Bot, interaction *discordgo.InteractionCreate){
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
			fmt.Printf("Error with unbanning member from guild: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "<@"+interaction.ApplicationCommandData().Options[0].StringValue()+"> is unbanned", interaction.Interaction)
	},

	"mute": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		member, err := bot.Session.GuildMember(interaction.GuildID, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).ID)

		err = bot.Session.GuildMemberMute(interaction.GuildID, member.User.ID, !member.Mute)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error has occurred", interaction.Interaction, time.Second*5)
			fmt.Printf("Error with muting member from guild: %s\n", err)
			return
		}
		if !member.Mute {
			interactions.SendMessageInteraction(bot.Session, member.Mention()+" is muted", interaction.Interaction)
			return
		}
		interactions.SendMessageInteraction(bot.Session, member.Mention()+" is unmuted", interaction.Interaction)
	},

	"defan": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		member, err := bot.Session.GuildMember(interaction.GuildID, interaction.ApplicationCommandData().Options[0].UserValue(bot.Session).ID)

		err = bot.Session.GuildMemberDeafen(interaction.GuildID, member.User.ID, !member.Deaf)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error has occurred", interaction.Interaction, time.Second*5)
			fmt.Printf("Error with muting member from guild: %s\n", err)
			return
		}
		if !member.Deaf {
			interactions.SendMessageInteraction(bot.Session, member.Mention()+" is defaned", interaction.Interaction)
			return
		}
		interactions.SendMessageInteraction(bot.Session, member.Mention()+" is undefaned", interaction.Interaction)
	},

	"purge": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		count := int(interaction.ApplicationCommandData().Options[0].IntValue())
		if count > 100 {
			count = 100
		}

		msgs, err := bot.Session.ChannelMessages(interaction.ChannelID, count, "", "", "")
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Can't delete last messages", interaction.Interaction, time.Second*5)
			fmt.Printf("Error with getting last messages: %s\n", err)
		}
		msgIds := make([]string, 0)

		for _, msg := range msgs {
			if msg.ID != "" {
				msgIds = append(msgIds, msg.ID)
			}
		}

		err = bot.Session.ChannelMessagesBulkDelete(interaction.ChannelID, msgIds)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Can't delete last messages", interaction.Interaction, time.Second*5)
			fmt.Printf("Error with delteing last messages: %s\n", err)
		}
		interactions.SendAndDeleteInteraction(bot.Session, "Last "+strconv.Itoa(count)+" messages are deleted", interaction.Interaction, time.Second*5)
	},

	"create_role": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		roleName := interaction.ApplicationCommandData().Options[0].StringValue()
		color := interaction.ApplicationCommandData().Options[1].IntValue()
		hoist := interaction.ApplicationCommandData().Options[2].BoolValue()
		permissions := interaction.ApplicationCommandData().Options[3].IntValue()
		mention := interaction.ApplicationCommandData().Options[4].BoolValue()

		role, err := bot.Session.GuildRoleCreate(interaction.GuildID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Role can't be created.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when creating new role: %s\n", err)
			return
		}

		role, err = bot.Session.GuildRoleEdit(interaction.GuildID, role.ID, roleName, int(color), hoist, permissions, mention)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Role can't be created.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when creating new role: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "Role is created "+role.Mention(), interaction.Interaction)
	},

	"edit_role": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		if len(interaction.ApplicationCommandData().Options) < 2 {
			interactions.SendAndDeleteInteraction(bot.Session, "Please select option to edit", interaction.Interaction, time.Second*5)
		}

		role := interaction.ApplicationCommandData().Options[0].RoleValue(bot.Session, interaction.GuildID)

		var (
			name       = role.Name
			color      = role.Color
			hoist      = role.Hoist
			permission = role.Permissions
			mention    = role.Mentionable
		)

		for _, option := range interaction.ApplicationCommandData().Options {
			switch option.Name {
			case "name":
				name = option.StringValue()
				break
			case "color":
				color = int(option.IntValue())
				break
			case "hoist":
				hoist = option.BoolValue()
				break
			case "permission":
				permission = option.IntValue()
				break
			case "mention":
				mention = option.BoolValue()
				break
			}

		}

		role, err := bot.Session.GuildRoleEdit(interaction.GuildID, role.ID, name, int(color), hoist, permission, mention)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Role can't be edited.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when updating role: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "Role is updated "+role.Mention(), interaction.Interaction)
	},

	"delete_role": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		role := interaction.ApplicationCommandData().Options[0].RoleValue(bot.Session, interaction.GuildID)

		err := bot.Session.GuildRoleDelete(interaction.GuildID, role.ID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Role can't be deleted.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when deleting role: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "Role is deleted.", interaction.Interaction)
	},

	"add_role": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		user := interaction.ApplicationCommandData().Options[0].UserValue(bot.Session)
		role := interaction.ApplicationCommandData().Options[1].RoleValue(bot.Session, interaction.GuildID)
		member, err := bot.Session.GuildMember(interaction.GuildID, user.ID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error when getting user information.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when getting user information: %s\n", err)
			return
		}

		if findUserRole(member, role) {
			interactions.SendAndDeleteInteraction(bot.Session, user.Mention()+" already have this role", interaction.Interaction, time.Second*5)
		}

		err = bot.Session.GuildMemberRoleAdd(interaction.GuildID, user.ID, role.ID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error when setting user role.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when setting user role: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "Role is added to "+user.Mention(), interaction.Interaction)
	},

	"remove_role": func(bot *player.Bot, interaction *discordgo.InteractionCreate) {
		user := interaction.ApplicationCommandData().Options[0].UserValue(bot.Session)
		role := interaction.ApplicationCommandData().Options[1].RoleValue(bot.Session, interaction.GuildID)
		member, err := bot.Session.GuildMember(interaction.GuildID, user.ID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error when getting user information.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when getting user information: %s\n", err)
			return
		}

		if !findUserRole(member, role) {
			interactions.SendAndDeleteInteraction(bot.Session, user.Mention()+" doesn't have this role", interaction.Interaction, time.Second*5)
		}

		err = bot.Session.GuildMemberRoleRemove(interaction.GuildID, user.ID, role.ID)
		if err != nil {
			interactions.SendAndDeleteInteraction(bot.Session, "Error when setting user role.", interaction.Interaction, time.Second*5)
			fmt.Printf("Error when setting user role: %s\n", err)
			return
		}
		interactions.SendMessageInteraction(bot.Session, "Role is removed from "+user.Mention(), interaction.Interaction)
	},
}
