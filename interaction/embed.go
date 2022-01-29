package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// SendEmbedInteraction Sends an embed as response to an interaction
func SendEmbedInteraction(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction) {
	sliceEmbed := []*discordgo.MessageEmbed{embed}
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Embeds: sliceEmbed},
	})
	if err != nil {
		fmt.Printf("InteractionRespond failed: %s\n", err)
		return
	}
}

// SendAndDeleteEmbedInteraction Sends and delete after the specified wait time as response to an interaction
func SendAndDeleteEmbedInteraction(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction, wait time.Duration) {
	SendEmbedInteraction(session, embed, interaction)

	time.Sleep(wait)

	err := session.InteractionResponseDelete(session.State.User.ID, interaction)
	if err != nil {
		fmt.Printf("InteractionResponseDelete failed: %s\n", err)
		return
	}
}

// ModifyEmbedInteraction Modify an already sent interaction
func ModifyEmbedInteraction(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction) {
	sliceEmbed := []*discordgo.MessageEmbed{embed}
	_, err := session.InteractionResponseEdit(session.State.User.ID, interaction, &discordgo.WebhookEdit{Embeds: sliceEmbed})
	if err != nil {
		fmt.Printf("InteractionResponseEdit failed: %session\n", err)
		return
	}
}

// ModifyAndDeleteEmbedInteraction Modify an already sent interaction and deletes it after the specified wait time
func ModifyAndDeleteEmbedInteraction(session *discordgo.Session, embed *discordgo.MessageEmbed, interaction *discordgo.Interaction, wait time.Duration) {
	ModifyEmbedInteraction(session, embed, interaction)

	time.Sleep(wait)

	err := session.InteractionResponseDelete(session.State.User.ID, interaction)
	if err != nil {
		fmt.Printf("InteractionResponseDelete failed: %session\n", err)
		return
	}
}
