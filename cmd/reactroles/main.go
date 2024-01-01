package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/reactroles/internal/dgclient"
	"github.com/zaptross/reactroles/internal/pgdb"
)

type AppSettings struct {
	ChatCommands  bool `default:"false"`
	SlashCommands bool `default:"true"`
}

func main() {
	var postgresConfig pgdb.PostgresDbParams
	pgErr := envconfig.Process("postgres", &postgresConfig)

	if pgErr != nil {
		log.Fatal(pgErr.Error())
	}

	db := pgdb.GetDatabase(postgresConfig)

	var discordConfig dgclient.DiscordGoClientParams

	dgErr := envconfig.Process("discord", &discordConfig)

	if dgErr != nil {
		log.Fatal(dgErr.Error())
	}

	discordConfig.DB = db
	bot := dgclient.GetClient(discordConfig)

	var appSettings AppSettings
	asErr := envconfig.Process("reactroles", &appSettings)

	if asErr != nil {
		log.Fatal(asErr.Error())
	}

	if appSettings.ChatCommands {
		// message commands
		bot.Session.AddHandler(bot.GetOnMessageHandler())
	}
	log.Println(commandsEnabled("chat", appSettings.ChatCommands))

	if appSettings.SlashCommands {
		// slash commands
		bot.Session.AddHandler(bot.GetOnInteractionHandler())
		_, err := bot.Session.ApplicationCommandCreate(discordConfig.AppID, "", bot.GetSlashCommand())

		if err != nil {
			log.Fatal(err.Error())
		}
	}
	log.Println(commandsEnabled("slash", appSettings.SlashCommands))

	// handle reactions
	bot.Session.AddHandler(bot.GetOnReactionAddHandler())
	bot.Session.AddHandler(bot.GetOnReactionRemoveHandler())

	bot.Connect()
	bot.Disconnect()
}

func commandsEnabled(t string, e bool) string {
	return fmt.Sprintf("[bot] %s commands enabled: %t", t, e)
}
