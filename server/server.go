package main

import (
	"fmt"
	"net"
	"networking/LAB_4/connection"
	"os"
	// "networking/LAB_4/constants"
	"encoding/json"
)

// Stores the number of active connections
var state []connection.ServerConn

// List of ports which can be provided to the client chatting applications
var portList = []string{":9000", ":9001", ":9002", ":9003", ":9004", ":9005", ":9006", ":9007"}

func main() {
	// For connecting to the port the helper server is located on
	port := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", port)
	terminateOnError(err)
	// Start the server and start listening on :8000
	listener, err := net.ListenTCP("tcp", tcpAddr)
	terminateOnError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			break
		}
		// Spawn a new goroutine for new Connections
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	err := demultiplexConn(conn)
	if err != nil {
		fmt.Println(err.Error())
		conn.Close()
	}
}

/*
This function figures what kind of request it is
There are 3 kinds of request
1. "A" : "For adding a new user to the chat and assign a new port number to the user and update the user state list"
2. "R" : For removing a user from the chat group releasing the assigned port and update the user list
3. "G" : For getting the list of current connected users
*/
func demultiplexConn(conn net.Conn) error {
	buffer := make([]byte, 100)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}
	if string(buffer[0:1]) == "A" {
		name := string(buffer[1:n])
		details := connection.ServerConn{Name: name, AddressString: portList[0]}
		portList = portList[1:]
		state = append(state, details)
	} else if string(buffer[0:1]) == "R" {
		name := string(buffer[1:n])
		removeClientByName(name)
	} else if string(buffer[0:1]) == "G" {
		b, err := json.Marshal(state)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		_, err = conn.Write(b)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
Function which removes the client from the user list by name
*/
func removeClientByName(name string) {
	index := indexOf(connection.ServerConn{Name: name, AddressString: ""})
	if index != -1 {
		portList = append(portList, state[index].AddressString)
		state = append(state[0:index], state[index+1:]...)
	}
}

/*
Function which removes the client from the user list by address
*/
func removeClientByAddress(addr string) {
	index := indexOf(connection.ServerConn{Name: "", AddressString: addr})
	if index != -1 {
		state = append(state[0:index], state[index+1:]...)
	}
}

func terminateOnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(2)
	}
}

func printErrorAndContinue(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
func indexOf(detail connection.ServerConn) int {
	for index, v := range state {
		if v.Name == detail.Name || v.AddressString == detail.AddressString {
			return index
		}
	}
	return -1
}
