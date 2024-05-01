package main

import (
	// "ticket_reservation/src/client"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Client struct {
	httpClient *http.Client
	eventURL   string
	reserveURL string
}

type Event struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Date             string `json:"date"`
	TotalTickets     int    `json:"total_tickets"`
	AvailableTickets int    `json:"available_tickets"`
}

type Response struct {
	Events []*Event `json:"events"`
}

func NewClient(eventURL, reserveURL string) *Client {
	return &Client{
		httpClient: &http.Client{},
		eventURL:   eventURL,
		reserveURL: reserveURL,
	}
}

func StartClient() {
	log.Println("Client is running")
	// eventURL := "http://localhost:8080/events"
	// reserveURL := "http://localhost:8080/events/"
	// client := NewClient(eventURL, reserveURL)
	// GetInput(client)
}

func (c *Client) GetEventsHandler(args []string) {
	if len(args) != 0 {
		log.Println("Invalid argument for getEvents command")
		return
	}

	resp, err := c.httpClient.Get(c.eventURL)
	if err != nil {
		log.Println("Error getting events:", err)
		return
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	var response Response
	if err := json.Unmarshal(responseData, &response); err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Status code %d is unexpected", resp.StatusCode)
	}

	PrintListOfEvents(response)
}

func (c *Client) BookTicketsHandler(args []string) {
	if len(args) != 2 {
		fmt.Println("Invalid arguments for bookEvents command")
		return
	}

	eventID := args[0]
	numTickets, err := strconv.Atoi(args[1])
	if err != nil {
		log.Println("Invalid number of tickets:", err)
		return
	}

	ticketResponse, err := json.Marshal(map[string]interface{}{
		"event_id":   args[0],
		"ticket_ids": args[1],
	})

	url := fmt.Sprintf("%s%s/reserve?num_tickets=%d", c.reserveURL, eventID, numTickets)
	resp, err := c.httpClient.Post(url, "", bytes.NewBuffer(ticketResponse))
	if err != nil {
		log.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	PrintBookTickets(body)
}
