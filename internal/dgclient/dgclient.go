package dgclient

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type DiscordGoClientParams struct {
	Token            string
	RoleMessage      string
	RoleChannel      string
	RoleAddRoleID    string
	RoleRemoveRoleID string
	DB               *pgdb.ReactRolesDatabase `ignored:"true"`
}

type DiscordGoClient struct {
	Session          *discordgo.Session
	roleMessage      *discordgo.Message
	roleAddRoleID    string
	roleRemoveRoleID string
	db               *pgdb.ReactRolesDatabase
}

func GetClient(params DiscordGoClientParams) *DiscordGoClient {
	if params.RoleChannel == "" {
		log.Fatal("Role channel not set in client params")
	}

	dg, err := discordgo.New("Bot " + params.Token)
	if err != nil {
		log.Fatal(err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages)

	client := &DiscordGoClient{
		Session:          dg,
		roleAddRoleID:    params.RoleAddRoleID,
		roleRemoveRoleID: params.RoleRemoveRoleID,
		db:               params.DB,
	}

	version, err := getVersionMessageIfPossible()

	if err == nil {
		log.Printf("[dgclient] Version: %s\n", version)
	}

	if params.RoleChannel != "" && params.RoleMessage == "" {
		message, err := dg.ChannelMessageSend(params.RoleChannel, "Setting up role assignment message...")

		if err != nil {
			log.Fatal(err)
		}

		client.roleMessage = message

		log.Printf("[dgclient] Role message created: %s\n", message.ID)
	}

	if params.RoleMessage != "" {
		message, err := dg.ChannelMessage(params.RoleChannel, params.RoleMessage)

		if err != nil {
			log.Fatal(err)
		}

		client.roleMessage = message
		log.Println("[dgclient] Role message found...")
	}

	return client
}

func (d *DiscordGoClient) Connect() {
	err := d.Session.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[dgclient] Connected to Discord, updating role message...")
	d.updateRoleSelectorMessage()

	log.Println("[dgclient] Waiting for events...")

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (d *DiscordGoClient) Disconnect() {
	d.Session.Close()
	println("Disconnected from Discord, exiting...")
}
