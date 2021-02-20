package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"networking/LAB_4/connection"
	"os"
	"strings"
)

var clientList connection.JSONServerConn

func main() {
	// Get the name the client wants to be known over the network
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your client name(this name will reference you in the connection):\n")
	name, _ := reader.ReadString('\n')
	// Defers the execution of the remove client function when the main function returns
	defer removeClient(name)
	// Add client to the server active users list
	err := addClient(name)
	terminateOnError(err)

	// Get the list of active users and their ports from the server
	list, err := getList()
	terminateOnError(err)
	// Start listening to the port assigned by the server
	listeningPort := getPort(list, name)
	clienttcpAddr, err := net.ResolveTCPAddr("tcp", listeningPort)
	terminateOnError(err)
	listener, err := net.ListenTCP("tcp", clienttcpAddr)
	terminateOnError(err)
	// Whenever there is a connection request to this chat client it spawns a goroutine and goes of the main thread
	go listen(listener)
	for {
		
		fmt.Println("Send Message(Y)")
		fmt.Println("Exit(E)")
		choice, _ := reader.ReadString('\n')
		if choice == "Y\n" {
			list, _ := getList()
			fmt.Printf("People online :::  %d\n", len(list))

			for _, v := range list {
				if v.Name != name {
					fmt.Printf("- %s\n", v.Name)
				}
			}
			message, _ := reader.ReadString('\n')
			index := strings.Index(message, "/")
			runningPort := getPort(list, message[0:index])
			runningTCPAddr, err := net.ResolveTCPAddr("tcp", runningPort)
			conn, err := net.DialTCP("tcp", nil, runningTCPAddr)
			printErrorAndContinue(err)
			conn.Write([]byte(message[index+1:]))
		} else if choice == "E\n" {
			break
		}
	}

	// conn.
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

func getList() ([]connection.JSONServerConn, error) {
	// Connecting to the port
	connectingPort := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", connectingPort)
	terminateOnError(err)
	var list []connection.JSONServerConn
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return list, err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("G"))
	if err != nil {
		return list, err
	}
	enc := json.NewDecoder(conn)
	err = enc.Decode(&list)
	if err != nil {
		return list, err
	}
	return list, nil
}

func getPort(connections []connection.JSONServerConn, name string) string {
	for _, v := range connections {
		if strings.Trim(v.Name, "\n") == name || v.Name == name {
			return v.AddressString
		}
	}
	return ""
}

func listen(listener *net.TCPListener) {
	buffer := make([]byte, 100)
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		n, _ := conn.Read(buffer)
		fmt.Printf(string(buffer[0:n]))
		conn.Close()
	}
}

func getPortFromString(str string) string {
	res := strings.Split(str, ":")
	return res[1]
}

func addClient(name string) error {
	// Connecting to the port
	connectingPort := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", connectingPort)
	terminateOnError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer conn.Close()
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte("A" + name))
	if err != nil {
		return err
	}
	return nil
}

func removeClient(name string) error {
	// Connecting to the port
	connectingPort := ":8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp", connectingPort)
	terminateOnError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte("R" + name))
	if err != nil {
		return err
	}
	return nil
}
