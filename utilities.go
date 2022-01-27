package main

import (
	"context"
	"fmt"
	"github.com/DisgoOrg/disgolink/dgolink"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/bwmarrin/discordgo"
	"net/url"
	"os"
	"strconv"
	"time"
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

func sendEmbed(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction) {
	sliceEmbed := []*discordgo.MessageEmbed{embed}
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{Type: discordgo.InteractionResponseChannelMessageWithSource, Data: &discordgo.InteractionResponseData{Embeds: sliceEmbed}})
	if err != nil {
		fmt.Printf("InteractionRespond failed: %s\n", err)
		return
	}
}

// Sends and delete after three second an embed in a given channel
func sendAndDeleteEmbedInteraction(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction, wait time.Duration) {
	sendEmbed(session, embed, interaction)

	time.Sleep(wait)

	err := session.InteractionResponseDelete(session.State.User.ID, interaction)
	if err != nil {
		fmt.Printf("InteractionResponseDelete failed: %s\n", err)
		return
	}
}

func validURL(value string) bool {
	_, err := url.ParseRequestURI(value)
	return err == nil
}

func play(s *discordgo.Session, link *dgolink.Link, guildID, voiceChannelID, channelID string, track lavalink.AudioTrack) {
	if err := s.ChannelVoiceJoinManual(guildID, voiceChannelID, false, false); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while joining voice channel: "+err.Error())
		return
	}
	if err := link.Player(snowflake.Snowflake(guildID)).Play(track); err != nil {
		_, _ = s.ChannelMessageSend(channelID, "error while playing track: "+err.Error())
		return
	}
	_, _ = s.ChannelMessageSend(channelID, "Playing: "+track.Info().Title())
}

func (b *Bot) registerNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	_, err := b.Link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:     "test",
		Host:     os.Getenv("LAVALINK_HOST"),
		Port:     os.Getenv("LAVALINK_PORT"),
		Password: os.Getenv("LAVALINK_PASSWORD"),
		Secure:   secure,
	})
	if err != nil {
		fmt.Printf("Error with registreing lavalink%s\n", err)
	}
}
