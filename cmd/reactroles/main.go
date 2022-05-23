package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		panic(err)
	}
	println("Bot is now running.")

	dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		println(fmt.Sprintf("%s: %s", m.Author.Username, m.Content))

		if m.Content == "!ping" {
			reactErr := s.MessageReactionAdd(m.ChannelID, m.ID, "🏓")

			if reactErr != nil {
				println(reactErr.Error())
			}
		}
	})

	strconv.ParseInt(os.Getenv("DISCORD_BOT_INTENT"), 10, 64)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuildMessages)

	dg.Open()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Open()
}
