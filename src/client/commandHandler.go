package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

func HelpHandler() {
	log.Println("Available commands:")
	log.Println("create [event name] [event date] [total tickets]")
	log.Println("get")
	log.Println("book [eventID] [number of tickets]")
	log.Println("exit")
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

func PrintCreateEvent(body []byte) {
	var event map[string]interface{}
	err := json.Unmarshal(body, &event)
	if err != nil {
		log.Println("Error unmarshalling response: " + err.Error())
		return
	}
	fmt.Println("Event created successfully")
	fmt.Println("ID:", event["ID"])
	fmt.Println("Name:", event["Name"])
	fmt.Println("Date:", event["Date"])
	fmt.Println("Total Tickets:", event["TotalTickets"])
	fmt.Println("Available Tickets:", event["AvailableTickets"])
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
		case "get":
			Client.GetEventsHandler(args)
		case "book":
			Client.BookTicketsHandler(args)
		case "create":
			Client.CreateEventHandler(args)
		case "help":
			HelpHandler()
		case "exit":
			ExitHandler()
		default:
			log.Println("Invalid command. Type 'help' to view available commands")
		}
	}
}
