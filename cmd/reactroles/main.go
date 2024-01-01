package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/zaptross/reactroles/internal/dgclient"
	"github.com/zaptross/reactroles/internal/pgdb"
)

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

	// message commands
	bot.Session.AddHandler(bot.GetOnMessageHandler())

	// reactions
	bot.Session.AddHandler(bot.GetOnReactionAddHandler())
	bot.Session.AddHandler(bot.GetOnReactionRemoveHandler())

	// slash commands
	bot.Session.AddHandler(bot.GetOnInteractionHandler())
	_, err := bot.Session.ApplicationCommandCreate(discordConfig.AppID, "", bot.GetSlashCommand())

	if err != nil {
		log.Fatal(err.Error())
	}

	bot.Connect()

	bot.Disconnect()
}
