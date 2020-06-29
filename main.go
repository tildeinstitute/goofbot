package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"time"

	"github.com/lrstanley/girc"
)

// ZW is a zero-width character
const ZW = string(0x200b)

// Generic in-chat error message
const errMsg = "I'm broke! Fix me! Check my logs!"

// Conf holds all the config info
type Conf struct {
	Owner  string `json:"owner"`
	Chan   string `json:"chan"`
	Server string `json:"server"`
	Port   int    `json:"port"`
	Nick   string `json:"nick"`
	Pass   string `json:"pass"`
	User   string `json:"user"`
	Name   string `json:"name"`
	SSL    bool   `json:"ssl"`
}

func main() {
	// check for config file specified by command line flag -c
	jsonLocation := flag.String("c", "config.json", "Path to config file in JSON format")
	jsonLocationLong := flag.String("config", "config.json", "Same as -c")

	// spit out config file structure if requested
	jsonFormat := flag.Bool("j", false, "Describes JSON config file fields")
	jsonFormatLong := flag.Bool("json", false, "Same as -j")

	flag.Parse()
	if *jsonFormat || *jsonFormatLong {
		fmt.Println(`Here is the format for the JSON config file:
            {
                "owner": "YourNickHere",
                "chan": "#bots",
                "server": "irc.tilde.chat",
                "port": 6697,
                "nick": "goofbot",
                "pass": "",
                "user": "goofbot",
                "name": "Goofus McBotus",
                "ssl": true
            }`)
		os.Exit(0)
	}

	// if the extended switch isn't the default value
	// and if the extended switch isn't empty
	if *jsonLocationLong != "config.json" && *jsonLocationLong != "" {
		*jsonLocation = *jsonLocationLong
	}

	// read the config file into a byte array
	jsonconf, err := ioutil.ReadFile(*jsonLocation)
	if err != nil {
		log.Fatalf("Error loading config: %v", err.Error())
	}

	// unmarshal the json byte array into struct conf
	var conf Conf
	err = json.Unmarshal(jsonconf, &conf)
	if err != nil {
		log.Fatalf("Error parsing config: %v", err.Error())
	}

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

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		// authenticate with nickserv if pass is set in config file
		if conf.Pass != "" {
			c.Cmd.Message("nickserv", "identify "+conf.Pass)
			time.Sleep(500 * time.Millisecond)
		}
		// join initial channel specified in config.json
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

		// when requested by owner, join channel specified
		if strings.HasPrefix(e.Last(), "!join") && e.Source.Name == conf.Owner {
			dest := strings.Split(e.Params[1], " ")
			if len(dest) < 2 {
				c.Cmd.Reply(e, "You must specify at least one channel to join")
				return
			}

			c.Cmd.Reply(e, "Right away, cap'n!")
			time.Sleep(100 * time.Millisecond)

			for _, channel := range dest[1:] {
				if !strings.HasPrefix(channel, "#") {
					continue
				}
				c.Cmd.Join(channel)
			}

			return
		}

		// respond with uptime / load
		if strings.HasPrefix(e.Last(), "!uptime") {
			uptime := exec.Command("/usr/bin/uptime")
			var out bytes.Buffer
			uptime.Stdout = &out
			err := uptime.Run()
			if err != nil {
				log.Printf("!uptime error: %v", err.Error())
				c.Cmd.Reply(e, errMsg)
				return
			}
			c.Cmd.Reply(e, out.String())
			return
		}

		if strings.HasPrefix(e.Last(), "!users") {
			who := exec.Command("/usr/local/bin/showwhoison", "")
			var bytestream bytes.Buffer
			who.Stdout = &bytestream
			err := who.Run()
			if err != nil {
				log.Printf("!users error: %v", err.Error())
				c.Cmd.Reply(e, errMsg)
				return
			}

			split := strings.Split(bytestream.String(), "\n")
			var out bytes.Buffer

			for i := 1; i < len(split); i++ {
				if split[i] == "" || len(split[i]) < 2 {
					continue
				}

				c := fmt.Sprintf("%s%s%s", split[i][:1], ZW, split[i][1:])
				out.WriteString(c + " ")
			}

			c.Cmd.Reply(e, out.String())
			return
		}

		// number of total human users on the server.
		// only active sessions.
		if strings.HasPrefix(e.Last(), "!totalusers") {
			userdirs, err := ioutil.ReadDir("/home")
			if err != nil {
				log.Printf("!totalusers error: %v", err.Error())
				c.Cmd.Reply(e, errMsg)
				return
			}

			msg := fmt.Sprintf("%v user accounts on ~institute", len(userdirs))
			c.Cmd.Reply(e, msg)
			return
		}

		// hmmm.
		if strings.Contains(e.Last(), "rain drop") {
			c.Cmd.Reply(e, "drop top")
		}

		if strings.HasPrefix(e.Last(), "!admin") {
			// gotify.sh contains a preconstructed curl request that
			// uses the gotify api to send a notification to admins
			gotify := exec.Command("./gotify.sh")
			err := gotify.Run()
			if err != nil {
				log.Printf("!admin error: %v", err.Error())
				c.Cmd.Reply(e, errMsg)
				return
			}

			c.Cmd.Reply(e, "The admins have been notified that you need their assistance!")
			return
		}
	})

	// die if there's a connection error
	if err := client.Connect(); err != nil {
		log.Fatalf("an error occurred while attempting to connect to %s: %s", client.Server(), err)
	}

	// sigint handling
	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt)
	go func() {
		<-ctrlc
		client.Close()
		os.Exit(1)
	}()
}
