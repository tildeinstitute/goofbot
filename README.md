# goofbot

An IRC bot I'm working on to practice Go.

## Features

External config file in JSON, optionally set via flag '-c' or '--config'  
Gives example config file with flag '-j' or '--json'  
Standard command/response structure  
* !totalusers - reports number of registered users  
* !users - reports logged in users  
* !uptime - reports uptime and load      
* !admin - Interacts with Gotify API to send a push notification to admins  
* !join #channel - Directs bot to join #channel  

Can define a bot owner for certain commands  
Able to identify with services  

## TODO

Externalize basic command/responses   
