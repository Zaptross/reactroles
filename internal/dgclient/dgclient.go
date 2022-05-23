package dgclient

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type DiscordGoClientParams struct {
	Token       string
	RoleMessage string
	RoleChannel string
	db          *pgdb.ReactRolesDatabase
}

type DiscordGoClient struct {
	Session     *discordgo.Session
	roleMessage *discordgo.Message
	db          *pgdb.ReactRolesDatabase
}

func GetClient(params DiscordGoClientParams) *DiscordGoClient {
	if params.RoleChannel == "" {
		panic("Role channel not set in client params")
	}

	dg, err := discordgo.New("Bot " + params.Token)
	if err != nil {
		panic(err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages)

	client := &DiscordGoClient{
		Session: dg,
		db:      params.db,
	}

	if params.RoleChannel != "" && params.RoleMessage == "" {
		message, err := dg.ChannelMessageSend(params.RoleChannel, "Setting up role assignment message...")

		if err != nil {
			panic(err)
		}

		client.roleMessage = message

		fmt.Printf("\nRole message created: %s\n\n", message.ID)
	}

	if params.RoleMessage != "" {
		message, err := dg.ChannelMessage(params.RoleChannel, params.RoleMessage)

		if err != nil {
			panic(err)
		}

		client.roleMessage = message
		println("Role message found...")
	}

	return client
}

func (d *DiscordGoClient) Connect() {
	err := d.Session.Open()
	if err != nil {
		panic(err)
	}

	println("Connected to Discord, waiting for events...")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (d *DiscordGoClient) Disconnect() {
	d.Session.Close()
	println("Disconnected from Discord, exiting...")
}
