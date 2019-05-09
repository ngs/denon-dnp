package main

// Simple program to send commands to Denon AVR and get their result

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

var (
	gr      chan string // global receiving channel for information coming from the AVR
	conn    net.Conn    // global network connection to the AVR
	debug   = true
	verbose = true
)

func sendCmd(cmd string) {
	cmd = strings.ToUpper(cmd)
	if verbose {
		fmt.Println("Sending: ", cmd)
	}

	cmd = cmd + "\r"
	fmt.Fprintf(conn, cmd)
	time.Sleep(210 * time.Millisecond)

}

func receiver() {
	if debug {
		fmt.Println("Receiver started")
	}

	status, err := bufio.NewReader(conn).ReadString('\r')
	gr <- status
	for err == nil { // There must be more information..keep reading.
		status, err = bufio.NewReader(conn).ReadString('\r')
		gr <- status
	}

	if debug {
		fmt.Println("Receiver has stopped.")
	}

}

func printReceived() {
	if debug {
		fmt.Println("printReceived started")
	}
	for recievedMsg := range gr {
		if recievedMsg != "" {
			fmt.Println("received: ", recievedMsg)
		} else {
			if verbose {
				fmt.Println("Received no result.")
			}
		}
	}

	fmt.Println("Done printing received channel.")
}

func init() {
	if debug {
		fmt.Println("Initilizing global channels.")
	}
	gr = make(chan string)

	if debug {
		fmt.Print("Connecting..")
	}
	lconn, err := net.Dial("tcp", "192.168.1.9:23")

	if err != nil {
		fmt.Println("Connection failed")
		os.Exit(1)
	}
	if debug {
		fmt.Println("connected.")
	}
	conn = lconn // Probably a better pattern for this..

	go receiver()
	go printReceived()

}

func cursor(cmd string) {
	sendCmd(fmt.Sprintf("NS9%s", cmd))
	time.Sleep(1000 * time.Millisecond)
	sendCmd("NSA")
}

func main() {

	// cmd_seq := []string{"MU?", "MUOFF", "MU?","MUON","MU?"}
	cmdSeq := os.Args[1:]

	defer close(gr)
	defer conn.Close()

	for _, cmd := range cmdSeq {
		switch cmd {
		case "radio":
			sendCmd("SITUNER")
		case "p1":
			sendCmd("NSP1")
		case "p2":
			sendCmd("NSP2")
		case "p3":
			sendCmd("NSP3")
		case "usb":
			sendCmd("SIUSB")
		case "info":
			sendCmd("NSA")
		case "next":
			cursor("X")
		case "prev":
			cursor("Y")
		case "up":
			cursor("0")
		case "down":
			cursor("1")
		case "left":
			cursor("2")
		case "right":
			cursor("3")
		case "enter":
			cursor("4")
		case "off":
			sendCmd("PWSTANDBY")
		case "on":
			sendCmd("PWON")
		default:
			sendCmd(cmd)
		}

		// Do we need to wait between sending commands?
		// Probably not, but makes it easier to see whats going on during dev
		if debug {
			time.Sleep(1000 * time.Millisecond)
		}

	}
}
