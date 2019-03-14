package main

import (
	"crypto/tls"
	"log"
	"net"

	"gopkg.in/irc.v3"
)

func setConfs() (bool, irc.ClientConfig, tls.Config) {
	//set this to false if you need
	useTLS := true

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
				c.Write("JOIN #institute") //initial channel join
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

	// set up the tls params for the connection
	// see: https://golang.org/pkg/crypto/tls/#Config
	tlsconfig := tls.Config{
		InsecureSkipVerify:       false, //set to true if you want to be dumb
		RootCAs:                  nil,   //use the OS's root CAs
		PreferServerCipherSuites: true,  //use the server's cipher list
	}
	return useTLS, config, tlsconfig
}

func main() {
	useTLS, config, tlsconfig := setConfs()

	switch useTLS {
	case true:
		conn, err := tls.Dial("tcp", "irc.tilde.chat:6697", &tlsconfig)
	case false:
		conn, err := net.Dial("tcp", "irc.tilde.chat:6667")
	}
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
