package dgclient

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/zaptross/reactroles/internal/pgdb"
	"github.com/zaptross/reactroles/internal/utils"
)

type DiscordGoClientParams struct {
	Token string
	AppID string
	DB    *pgdb.ReactRolesDatabase `ignored:"true"`
}

type DiscordGoClient struct {
	Session *discordgo.Session
	db      *pgdb.ReactRolesDatabase
}

func GetClient(params DiscordGoClientParams) *DiscordGoClient {
	dg, err := discordgo.New("Bot " + params.Token)
	if err != nil {
		log.Fatal(err)
	}

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages)

	client := &DiscordGoClient{
		Session: dg,
		db:      params.DB,
	}

	s, c := utils.GetVersionRaw()
	version := fmt.Sprintf("%s (%s)", s, c)

	if err == nil {
		log.Printf("[dgclient] Version: %s\n", version)
	}

	return client
}

func (d *DiscordGoClient) Connect() {
	err := d.Session.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("[dgclient] Connected to Discord, updating role messages...")
	d.updateAllRoleSelectorMessages()

	log.Println("[dgclient] Waiting for events...")

	s, c := utils.GetVersionRaw()
	slog.Info("reactroles started successfully", "semantic", s, "commit", c)
	d.Session.UpdateCustomStatus(fmt.Sprintf("Version: %s (%s)", s, c))

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (d *DiscordGoClient) Disconnect() {
	d.Session.Close()
	println("Disconnected from Discord, exiting...")
}
