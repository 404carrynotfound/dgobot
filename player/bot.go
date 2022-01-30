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
	Session        *discordgo.Session
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

// Play Start a new session for lava player and returns a message
func (b *Bot) Play(voiceChannelID string, interaction *discordgo.Interaction, playlist string, tracks ...lavalink.AudioTrack) {
	guildID := interaction.GuildID
	waitTime := time.Second * 10

	err := b.Session.ChannelVoiceJoinManual(guildID, voiceChannelID, false, true)
	if err != nil {
		interactions.SendAndDeleteInteraction(b.Session, "Error while joining voice channel: "+err.Error(), interaction, waitTime)
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
			interactions.SendAndDeleteInteraction(b.Session, "Error while playing track: "+err.Error(), interaction, waitTime)
			return
		}
		if playlist == "" {
			interactions.SendMessageInteraction(b.Session, "Playing: "+track.Info().Title, interaction)
		} else {
			interactions.SendMessageInteraction(b.Session, "Playing: "+playlist, interaction)
		}
		manager.playing = true
	} else {
		if playlist == "" {
			interactions.SendMessageInteraction(b.Session, "Added to queue: "+tracks[0].Info().Title, interaction)
		} else {
			interactions.SendMessageInteraction(b.Session, "Added to queue: "+playlist, interaction)
		}
	}

}

// Skip skips current song in lava player
func (b *Bot) Skip(interaction *discordgo.Interaction) {
	guildID := interaction.GuildID
	waitTime := time.Second * 10

	manager, ok := b.PlayerManagers[guildID]
	if !ok {
		return
	}

	if len(manager.Queue) == 0 {
		b.Stop(interaction)
	}

	if manager.playing {
		track := manager.PopQueue()
		if err := manager.Player.Play(track); err != nil {
			interactions.SendAndDeleteInteraction(b.Session, "Error while playing track: "+err.Error(), interaction, waitTime)
			return
		}
		interactions.SendMessageInteraction(b.Session, "Playing: "+track.Info().Title, interaction)
	}
}

// Stop stopping current lava player session, clears current queue and bot leave current voice channel
func (b *Bot) Stop(interaction *discordgo.Interaction) {
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

	err = b.Session.ChannelVoiceJoinManual(guildID, "", false, false)
	if err != nil {
		fmt.Printf("Error when leaving channel: %s\n", err)
		return
	}

	interactions.SendMessageInteraction(b.Session, "Queue is cleared", interaction)
}

// Current sends a message with current song in lava player
func (b *Bot) Current(interaction *discordgo.Interaction) {
	track := b.Link.Player(snowflake.Snowflake(interaction.GuildID)).Track().Info()
	interactions.SendMessageInteraction(b.Session, track.Author+" - "+track.Title, interaction)
}

// Pause pauses current lava player session
func (b *Bot) Pause(interaction *discordgo.Interaction) {
	guildSnowflake := snowflake.Snowflake(interaction.GuildID)
	if b.Link.Player(guildSnowflake).Paused() {
		interactions.SendAndDeleteInteraction(b.Session, "Player is already paused", interaction, time.Second*10)
		return
	}
	interactions.SendMessageInteraction(b.Session, "Player is paused", interaction)
	_ = b.Link.Player(guildSnowflake).Pause(true)
}

// Resume resumes current lava player session
func (b *Bot) Resume(interaction *discordgo.Interaction) {
	guildSnowflake := snowflake.Snowflake(interaction.GuildID)
	if !b.Link.Player(guildSnowflake).Paused() {
		interactions.SendAndDeleteInteraction(b.Session, "Player is already resumed", interaction, time.Second*10)
		return
	}
	interactions.SendMessageInteraction(b.Session, "Player is resumed", interaction)
	_ = b.Link.Player(guildSnowflake).Pause(false)
}
