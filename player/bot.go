package player

import (
	"context"
	interactions "dgobot/interaction"
	"fmt"
	"github.com/DisgoOrg/disgolink/dgolink"
	"github.com/DisgoOrg/disgolink/lavalink"
	"github.com/DisgoOrg/snowflake"
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
	"time"
)

type Bot struct {
	Link           *dgolink.Link
	PlayerManagers map[string]*Manager
}

func (b *Bot) RegisterNodes() {
	secure, _ := strconv.ParseBool(os.Getenv("LAVALINK_SECURE"))
	_, err := b.Link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:     "test",
		Host:     os.Getenv("LAVALINK_HOST"),
		Port:     os.Getenv("LAVALINK_PORT"),
		Password: os.Getenv("LAVALINK_PASSWORD"),
		Secure:   secure,
	})
	if err != nil {
		fmt.Printf("Error with registreing lavalink: %s\n", err)
	}
}

func (b *Bot) Play(session *discordgo.Session, voiceChannelID string, interaction *discordgo.Interaction, playlist string, tracks ...lavalink.AudioTrack) {
	guildID := interaction.GuildID
	waitTime := time.Second * 10

	err := session.ChannelVoiceJoinManual(guildID, voiceChannelID, false, true)
	if err != nil {
		interactions.SendAndDeleteInteraction(session, "Error while joining voice channel: "+err.Error(), interaction, waitTime)
		return
	}

	manager, ok := b.PlayerManagers[guildID]
	if !ok {
		manager = &Manager{
			Player: b.Link.Player(snowflake.Snowflake(guildID)),
		}
		b.PlayerManagers[guildID] = manager
		manager.Player.AddListener(manager)
	}

	manager.AddQueue(tracks...)

	if !manager.playing {
		track := manager.PopQueue()
		if err := manager.Player.Play(track); err != nil {
			interactions.SendAndDeleteInteraction(session, "Error while playing track: "+err.Error(), interaction, waitTime)
			return
		}
		if playlist == "" {
			interactions.SendMessageInteraction(session, "Playing: "+track.Info().Title, interaction)
		} else {
			interactions.SendMessageInteraction(session, "Playing: "+playlist, interaction)
		}
		manager.playing = true
	} else {
		if playlist == "" {
			interactions.SendMessageInteraction(session, "Playing: "+tracks[0].Info().Title, interaction)
		} else {
			interactions.SendMessageInteraction(session, "Playing: "+playlist, interaction)
		}
	}

}

func (b *Bot) Skip(session *discordgo.Session, interaction *discordgo.Interaction) {
	guildID := interaction.GuildID
	waitTime := time.Second * 10

	manager, ok := b.PlayerManagers[guildID]
	if !ok {
		return
	}

	if len(manager.Queue) == 0 {
		b.Stop(session, interaction)
	}

	if manager.playing {
		track := manager.PopQueue()
		if err := manager.Player.Play(track); err != nil {
			interactions.SendAndDeleteInteraction(session, "Error while playing track: "+err.Error(), interaction, waitTime)
			return
		}
		interactions.SendMessageInteraction(session, "Playing: "+track.Info().Title, interaction)
	}
}

func (b *Bot) Stop(session *discordgo.Session, interaction *discordgo.Interaction) {
	guildID := interaction.GuildID

	manager, ok := b.PlayerManagers[guildID]
	if !ok {
		return
	}

	manager.EmptyQueue()
	manager.playing = false
	err := manager.Player.Stop()
	if err != nil {
		fmt.Printf("Error when stopping player: %s\n", err)
		return
	}

	err = session.ChannelVoiceJoinManual(guildID, "", false, false)
	if err != nil {
		fmt.Printf("Error when leaving channel: %s\n", err)
		return
	}

	interactions.SendMessageInteraction(session, "Queue is cleared", interaction)
}

func (b *Bot) Current(session *discordgo.Session, interaction *discordgo.Interaction) {
	interactions.SendMessageInteraction(session, b.Link.Player(snowflake.Snowflake(interaction.GuildID)).Track().Info().Title, interaction)
}

func (b *Bot) Pause(session *discordgo.Session, interaction *discordgo.Interaction) {
	guildSnowflake := snowflake.Snowflake(interaction.GuildID)
	if b.Link.Player(guildSnowflake).Paused() {
		interactions.SendAndDeleteInteraction(session, "Player is already paused", interaction, time.Second*10)
		return
	}
	interactions.SendMessageInteraction(session, "Player is paused", interaction)
	_ = b.Link.Player(guildSnowflake).Pause(true)
}

func (b *Bot) Resume(session *discordgo.Session, interaction *discordgo.Interaction) {
	guildSnowflake := snowflake.Snowflake(interaction.GuildID)
	if !b.Link.Player(guildSnowflake).Paused() {
		interactions.SendAndDeleteInteraction(session, "Player is already resumed", interaction, time.Second*10)
		return
	}
	interactions.SendMessageInteraction(session, "Player is resumed", interaction)
	_ = b.Link.Player(guildSnowflake).Pause(false)
}
