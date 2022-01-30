package main

import (
	"github.com/bwmarrin/discordgo"
)

var Commands = []*discordgo.ApplicationCommand{
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
	{
		Name:        "mute",
		Description: "mute user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "mute user",
				Required:    true,
			},
		},
	},
	{
		Name:        "defan",
		Description: "deafan user",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "deafan user",
				Required:    true,
			},
		},
	},
	{
		Name:        "purge",
		Description: "delete last n messages",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "count",
				Description: "The number of the messages to be deleted. A maximum of 100 messages.",
				Required:    true,
			},
		},
	},
	{
		Name:        "create_role",
		Description: "Creates role in current guild.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "The name of the Role.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "color",
				Description: "The color of the role (decimal, not hex).",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "hoist",
				Description: "Whether to display the role's users separately.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "permission",
				Description: "The permissions for the role.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "mention",
				Description: "Whether this role is mentionable.",
				Required:    true,
			},
		},
	},
	{
		Name:        "edit_role",
		Description: "Edit role in current guild.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The name of the Role to be edited.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "The name of the Role.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "color",
				Description: "The color of the role (decimal, not hex).",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "hoist",
				Description: "Whether to display the role's users separately.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "permission",
				Description: "The permissions for the role.",
				Required:    false,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "mention",
				Description: "Whether this role is mentionable.",
				Required:    false,
			},
		},
	},
	{
		Name:        "delete_role",
		Description: "Deletes role in current guild.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "The name of the Role to be deleted.",
				Required:    true,
			},
		},
	},
	{
		Name:        "add_role",
		Description: "Adds role to user.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The name of the user to be added role.",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionRole,
				Name:        "role",
				Description: "Role to be assigned to the user",
				Required:    true,
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
	// TODO: GuildMemberMute
}
