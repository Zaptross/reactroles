package dgclient

import (
	"github.com/bwmarrin/discordgo"
)

func (client *DiscordGoClient) GetOnInteractionHandler() func(*discordgo.Session, *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: ":gear: Working...",
				Flags:   uint64(discordgo.MessageFlagsEphemeral),
			},
		})

		defer client.updateRoleSelectorMessage()

		if i.Member.User.ID == s.State.User.ID {
			return
		}

		switch i.ApplicationCommandData().Options[0].Name {
		case Actions.Add:
			handleAddRoleSlashCommand(client, s, i)
		case Actions.Remove:
			handleRemoveRoleSlashCommand(client, s, i)
		case Actions.Update:
			handleUpdateRoleSlashCommand(client, s, i)
		}
	}
}

func (client *DiscordGoClient) GetSlashCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "role",
		Description: "Manage roles",
		Options: []*discordgo.ApplicationCommandOption{
			addRoleSlashCommand(),
			removeRoleSlashCommand(),
			updateRoleSlashCommand(),
		},
	}
}
