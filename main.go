package main

import (
	"log"
	"net"

	"gopkg.in/irc.v3"
)

func main() {

	// set all the necessities
	// also specifies the initial channel
	// to join
	config := irc.ClientConfig{
		Nick: "goofbot",
		Pass: "password",
		User: "goofbot",
		Name: "Goofus McBot",
		Handler: irc.HandlerFunc(func(c *irc.Client, m *irc.Message) {
			if m.Command == "001" {
				// 001 = welcome
				c.Write("JOIN #institute")
			} else if m.Command == "PRIVMSG" && c.FromChannel(m) {
				c.WriteMessage(&irc.Message{
					Command: "PRIVMSG",
					Params: []string{
						m.Params[0],
						m.Trailing(),
					},
				})
			}
		}),
	}

	conn, err := net.Dial("tcp", "irc.tilde.chat:6697")
	if err != nil {
		log.Fatalln(err)
	}

	// create the connection
	client := irc.NewClient(conn, config)
	err = client.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
