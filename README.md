# goofbot

An IRC bot I'm working on to practice Go.

## Features

External config file in JSON, optionally set via flag '-c'  
Gives example config file with flag '-j'  
Standard command/response structure
    !totalusers - reports number of registered users  
    !users - reports logged in users  
    !uptime - reports uptime and load      
Can define a bot owner for certain commands  
Able to identify with services  

## TODO

Externalize basic command/responses  
Add interaction with Gotify API to send notification to channel admin when needed  
