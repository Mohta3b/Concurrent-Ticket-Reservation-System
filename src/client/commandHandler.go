package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"encoding/json"
)

func HelpHandler() {
	log.Println("Available commands:")
	log.Println("getEvents - Get a list of all events")
	log.Println("bookEvents [eventID] [number of tickets] - Book tickets for an event")
	log.Println("exit - Exit the program")
}

func ExitHandler() {
	log.Println("Exiting...")
	os.Exit(0)
}

func PrintListOfEvents(response Response) {
	for _, event := range response.Events {
		fmt.Println("Event ID:", event.ID)
		fmt.Println("Event Name:", event.Name)
		fmt.Println("Event Date:", event.Date)
		fmt.Println("Total Tickets:", event.TotalTickets)
		fmt.Println("Available Tickets:", event.AvailableTickets)
		fmt.Println()
	}
}

func PrintBookTickets(body []byte) {
	var ticketIDs []string
	err := json.Unmarshal(body, &ticketIDs)
	if err != nil {
		log.Println("Error unmarshalling response: " + err.Error())
		return
	}
	fmt.Println("Tickets booked successfully")
	fmt.Println("Ticket IDs:")
	for _, ticketID := range ticketIDs {
		fmt.Println(ticketID)
	}
}

func GetInput(Client *Client) {
	scanner := bufio.NewScanner(os.Stdin)
	log.Println("Please enter a command or type 'help' to view available commands: ")
	for scanner.Scan() {
		commandParts := strings.Split(scanner.Text(), " ")
		command := commandParts[0]
		var args []string
		if len(commandParts) > 1 {
			args = commandParts[1:]
		}

		switch command {
		case "getEvents":
			Client.GetEventsHandler(args)
		case "bookEvents":
			Client.BookTicketsHandler(args)
		case "help":
			HelpHandler()
		case "exit":
			ExitHandler()
		default:
			log.Println("Invalid command. Type 'help' to view available commands")
		}
	}
}
