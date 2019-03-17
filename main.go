package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lrstanley/girc"
)

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

//right now Conf.Pass isn't used, but i'm leaving it so the bot
//can have a registered nick
type Conf struct {
	Owner  string
	Chan   string
	Server string
	Port   int
	Nick   string
	Pass   string
	User   string
	Name   string
	SSL    bool
}

func main() {
	//read the file into a byte array
	jsonconf, err := ioutil.ReadFile("config.json")
	checkerr(err)

	// unmarshal the json byte array into struct conf
	var conf Conf
	err = json.Unmarshal(jsonconf, &conf)
	checkerr(err)

	// CLIENT CONFIG
	client := girc.New(girc.Config{
		Server: conf.Server,
		Port:   conf.Port,
		Nick:   conf.Nick,
		User:   conf.User,
		Name:   conf.Name,
		Out:    os.Stdout,
		SSL:    conf.SSL,
	})

	// join startup channel
	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		c.Cmd.Join(conf.Chan)
	})

	// basic command-response handler
	client.Handlers.Add(girc.PRIVMSG, func(c *girc.Client, e girc.Event) {
		if strings.HasPrefix(e.Last(), "!hello") {
			c.Cmd.ReplyTo(e, "henlo good fren!!")
			return
		}

		// check if the command was issued by a specific person before dying
		// i had to delve into girc/event.go to find e.Source.Name
		if strings.HasPrefix(e.Last(), "die, devil bird!") && e.Source.Name == conf.Owner {
			c.Cmd.Reply(e, "SQUAWWWWWK!!")
			time.Sleep(100 * time.Millisecond)
			c.Close()
			return
		}
		//another basic command/response. required information for the tildeverse
		if strings.HasPrefix(e.Last(), "!botlist") {
			c.Cmd.Reply(e, "Creator: ~a h r i m a n~ :: I'm the assistance bot for tilde.institute. Commands: !hello !join !uptime !users")
			return
		}
		// when requested by owner, join channel specified
		if strings.HasPrefix(e.Last(), "!join") && e.Source.Name == conf.Owner {
			dest := strings.Split(e.Params[1], " ")
			c.Cmd.Reply(e, "Right away, cap'n!")
			time.Sleep(100 * time.Millisecond)
			c.Cmd.Join(dest[1])
		}
		// respond with uptime / load
		if strings.HasPrefix(e.Last(), "!uptime") {
			uptime := exec.Command("uptime")
			var out bytes.Buffer
			uptime.Stdout = &out
			err := uptime.Run()
			if err != nil {
				log.Fatalln("Error while running 'uptime'")
			}
			c.Cmd.Reply(e, out.String())
			return
		}
		// respond with currently connected users
		// TODO: prepend names with _ to avoid pings in irc
		if strings.HasPrefix(e.Last(), "!users") {
			users := exec.Command("who", "-q")
			var out bytes.Buffer
			users.Stdout = &out
			err := users.Run()
			if err != nil {
				log.Fatalln("Error while running 'who -q'")
			}
			c.Cmd.Reply(e, out.String())
			return
		}
		// number of total users
		// bot dies when i run this
		// TODO: defuckulate this command
		if strings.HasPrefix(e.Last(), "!totalusers") {
			users := exec.Command("ls /home | wc -w")
			var out bytes.Buffer
			users.Stdout = &out
			err := users.Run()
			if err != nil {
				log.Fatalln("Error while running 'ls /home | wc -w'")
			}
			c.Cmd.Reply(e, out.String())
		}
	})

	// if err is not nothing, eg, if there's an error
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
