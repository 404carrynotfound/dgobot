package interactions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"time"
)

// SendMessageInteraction Sends a message in as response to an interaction
func SendMessageInteraction(session *discordgo.Session, message string, interaction *discordgo.Interaction) {
	err := session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{Content: message},
	})
	if err != nil {
		fmt.Printf("InteractionRespond failed: %s\n", err)
		return
	}
}

// SendAndDeleteInteraction Sends and delete after the specified wait time as response to an interaction
func SendAndDeleteInteraction(session *discordgo.Session, message string, interaction *discordgo.Interaction, wait time.Duration) {
	SendMessageInteraction(session, message, interaction)

	time.Sleep(wait)

	err := session.InteractionResponseDelete(session.State.User.ID, interaction)
	if err != nil {
		fmt.Printf("InteractionResponseDelete failed: %s\n", err)
		return
	}
}

// ModifyInteraction Modify an already sent interaction
func ModifyInteraction(session *discordgo.Session, message string, interaction *discordgo.Interaction) {
	_, err := session.InteractionResponseEdit(session.State.User.ID, interaction, &discordgo.WebhookEdit{Content: message})
	if err != nil {
		fmt.Printf("InteractionResponseEdit failed: %s\n", err)
		return
	}
}

// ModifyAndDeleteInteraction Modify an already sent interaction and deletes it after the specified wait time
func ModifyAndDeleteInteraction(session *discordgo.Session, message string, interaction *discordgo.Interaction, wait time.Duration) {
	ModifyInteraction(session, message, interaction)

	time.Sleep(wait)

	err := session.InteractionResponseDelete(session.State.User.ID, interaction)
	if err != nil {
		fmt.Printf("InteractionResponseDelete failed: %s\n", err)
		return
	}
}
