package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/lrstanley/girc"
)

func main() {
	client := girc.New(girc.Config{
		Server: "irc.tilde.chat",
		Port:   6697,
		Nick:   "goofbot",
		User:   "goofbot",
		Name:   "Goofus McBotus",
		Out:    os.Stdout,
		SSL:    true,
	})

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join("#goofbot")
	})

	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		if strings.HasPrefix(e.Last(), "!hello") {
			c.Cmd.ReplyTo(e, "henlo good fren!!")
			return
		}

		if strings.HasPrefix(e.Last(), "die, devil bird!") && e.Source.Name == "ahriman" {
			c.Cmd.Reply(e, "SQUAWWWWWK!!")
			time.Sleep(100 * time.Millisecond)
			c.Close()
			return
		}
		if strings.HasPrefix(e.Last(), "!botlist") {
			c.Cmd.Reply(e, "Creator: ahriman. I'm the assistance bot for ~institute. Commands: !hello, !stop")
			return
		}
	})

	if err := client.Connect(); err != nil {
		log.Fatalf("an error occurred while attempting to connect to %s: %s", client.Server(), err)
	}
	//TODO: figure out sigint handling
	//	ctrlc := make(chan os.Signal, 1)
	//	signal.Notify(ctrlc, os.Interrupt)
	//	go func() {
	//		<-ctrlc
	//		client.Close()
	//		os.Exit(1)
	//	}()
}
