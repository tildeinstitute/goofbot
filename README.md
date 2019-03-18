# goofbot

An IRC bot I'm working on to practice Go.

## Features

External config file in JSON, optionally set via flag '-c'  
Gives example config file with flag '-j'  
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
Sanitize output of !users so as to not ping them in IRC  
