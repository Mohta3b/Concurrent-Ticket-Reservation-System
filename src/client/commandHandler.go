package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func PrintListOfEvents(response []Response) {
	if len(response) == 0 {
		fmt.Println("No events found")
		return
	}
	fmt.Println("Events:")
	for _, event := range response {
		fmt.Println("ID:", event.ID)
		fmt.Println("Name:", event.Name)
		fmt.Println("Date:", event.Date)
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
	var response Response
	err := json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error unmarshalling response: " + err.Error())
		return
	}
	fmt.Println("Event created successfully")
	fmt.Println("ID:", response.ID)
	fmt.Println("Name:", response.Name)
	fmt.Println("Date:", response.Date)
	fmt.Println("Total Tickets:", response.TotalTickets)
	fmt.Println("Available Tickets:", response.AvailableTickets)
}

func ConnectToServer(Client *Client) error {
	statusCode := Client.GetHomePageHandler()
	if statusCode != http.StatusOK {
		log.Println("Error connecting to server. Exiting...")
		return fmt.Errorf("error connecting to server")
	}
	GetInput(Client)
	return nil
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
