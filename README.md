# goofbot ![gpl-3.0](https://img.shields.io/badge/license-GPLv3-brightgreen.svg "GPL v3") [![Go Report Card](https://goreportcard.com/badge/github.com/gbmor/goofbot)](https://goreportcard.com/report/github.com/gbmor/goofbot) [![GolangCI](https://img.shields.io/badge/golang%20ci-success-blue.svg)](https://golangci.com/r/github.com/gbmor/goofbot)


An IRC bot I'm working on to practice Go.

## Features

External config file in JSON, file name/path optionally set via flag '-c' or '--config'  
Gives example config file with flag '-j' or '--json'  
Standard command/response structure  
* !totalusers - reports number of registered users  
* !users - reports logged in users  
* !uptime - reports uptime and load      
* !admin - Interacts with Gotify API to send a push notification to admins  
* !join #channel - Directs bot to join #channel  

Can define a bot owner for certain commands  
Able to identify with services  
Handles ^C (SIGINT) gracefully  

## TODO

Externalize basic command/responses   
