package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/lrstanley/girc"
)

// function to grease error checking
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

//Conf ... right now Conf.Pass isn't used,
//but i'm leaving it so the bot
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
	//check for config file specified by command line flag -c
	jsonlocation := flag.String("c", "config.json", "Path to config file in JSON format")
	jsonlocationlong := flag.String("config", "config.json", "Same as -c")
	//spit out config file structure if requested
	jsonformat := flag.Bool("j", false, "Describes JSON config file fields")
	jsonformatlong := flag.Bool("json", false, "Same as -j")

	flag.Parse()
	if *jsonformat == true || *jsonformatlong == true {
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
	if *jsonlocationlong != "config.json" && *jsonlocationlong != "" {
		*jsonlocation = *jsonlocationlong
	}

	//read the config file into a byte array
	jsonconf, err := ioutil.ReadFile(*jsonlocation)
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

	client.Handlers.Add(girc.CONNECTED, func(c *girc.Client, e girc.Event) {
		//authenticate with nickserv if pass is set in config file
		if conf.Pass != "" {
			c.Cmd.Message("nickserv", "identify "+conf.Pass)
			time.Sleep(500 * time.Millisecond)
		}
		//join initial channel specified in config.json
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
			c.Cmd.Reply(e, "Creator: ~a h r i m a n~ :: I'm the assistance bot for tilde.institute. Commands: !hello !join !uptime !users !totalusers. If you need the assistance of an admin, issue !admin")
			return
		}
		// when requested by owner, join channel specified
		if strings.HasPrefix(e.Last(), "!join") && e.Source.Name == conf.Owner {
			dest := strings.Split(e.Params[1], " ")
			c.Cmd.Reply(e, "Right away, cap'n!")
			time.Sleep(100 * time.Millisecond)
			c.Cmd.Join(dest[1])
			return
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
			// execs: who -q | awk 'NR==1'
			// then saves the output to bytestream
			who := exec.Command("who", "-q")
			awk := exec.Command("awk", "NR==1")
			r, w := io.Pipe()
			who.Stdout = w
			awk.Stdin = r
			var bytestream bytes.Buffer
			awk.Stdout = &bytestream
			who.Start()
			awk.Start()
			who.Wait()
			w.Close()
			awk.Wait()
			r.Close()

			c.Cmd.Reply(e, bytestream.String())
			return
		}
		// number of total human users on the server
		if strings.HasPrefix(e.Last(), "!totalusers") {
			userdirs, err := ioutil.ReadDir("/home")
			checkerr(err)
			c.Cmd.Reply(e, strconv.Itoa(len(userdirs))+" user accounts on ~institute")
			return
		}
		if strings.HasPrefix(e.Last(), "!admin") {
			//gotify.sh contains a preconstructed curl request that
			//uses the gotify api to send a notification to admins
			gotify := exec.Command("./gotify.sh")
			err := gotify.Run()
			if err == nil {
				c.Cmd.Reply(e, "The admins have been notified that you need their assistance!")
			}
			checkerr(err)
			return
		}
	})

	// die if there's a connection error
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
