# GoTALK
## A CLI P2P chat application in go which supports 7(configurable) clients at a time
![](./record.gif)
### Design
    There are 2 main parts to the application
    - Client: The process on which the user can see other active users and send messages to send by prepending the message with the username/. For example
        ```
        username = Ritwiz
        Message = Ritwiz/Hello How are you ?
        ```
        The client first asks the user the name they want to be referenced as in the chat
        After entering the name the client gets the list of active clients, establishes a TCP listener server off the main thread, and asks the user whether he/she wants to 
        send a message or exit the application. 
        If the user wants to send a message they are given a list of active users on the server  and the user can select any of them and message them
        The message appears on the same terminal the user is working on
    
    - Server: Keeps a list of all the active clients
        The server has a list of 7 ports numbers (This can be manually increased). The server supports stateless connections. Once a client connects to the
        server it can execute 3 operations according to the ARG protocol [A custom naive protocol over TCP named just for the sake of it  : ) ]
        - A: Add the user to the list of active clients 
        The client sends a string preceded by character A and followed by the name of the client. For example if the name of the client is "RAM"
        then the client will send "ARAM" to the server. The server intercepts the message and adds the user to the active users list and assigns a free port to the user
        - R: Remove the user from the active list of users (when the user exits)
        The user will send the name of the client appended with the "R" to indicate the user wants to leave the chat. The port number is freed and is added to the free port list
        - G: Get the list of all connected clients and their port numbers

### Installation
1. Install [go](https://golang.org/doc/install)
2. Go the project directory and for compiling the server 
```
cd server
go build server.go
```
3. For compiling the client 
```
cd client
go build client.go
```

### Running
If you don't want to compile or install go the executable files are also attached with project for both the client and the server. The server can be started using
```
./server
```

For starting an instance of the client

```
./client
```
Multiple clients maybe started depending on the number of ports configured in the server