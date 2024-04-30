package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "sync"
    "strconv"
    "strings"
)

type Client struct {
    httpClient *http.Client
    eventURL  string
    reserveURL string
}

type Event struct {
    ID              string `json:"id"`
    Name            string `json:"name"`
    Date            string `json:"date"`
    TotalTickets    int    `json:"total_tickets"`
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

func createClient() *Client {
    eventURL := "http://localhost:8080/events"
    reserveURL := "http://localhost:8080/events/"

    return NewClient(eventURL, reserveURL)
}

func (c *Client) GetEvents() []*Event {
    resp, err := c.httpClient.Get(c.eventURL)
    if err != nil {
        fmt.Println("Error getting events:", err)
        return nil
    }
    defer resp.Body.Close()

    responseData, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Error reading response:", err)
        return nil
    }

    var response Response
    if err := json.Unmarshal(responseData, &response); err != nil {
        fmt.Println("Error decoding JSON:", err)
        return nil
    }

    for _, event := range response.Events {
        fmt.Println("Event ID:", event.ID)
        fmt.Println("Event Name:", event.Name)
        fmt.Println("Event Date:", event.Date)
        fmt.Println("Total Tickets:", event.TotalTickets)
        fmt.Println("Available Tickets:", event.AvailableTickets)
        fmt.Println()
    }

    return response.Events
}

func (c *Client) BookTickets(eventID string, numTickets int) ([]string, error) {
	
	resp, err := c.httpClient.Post(fmt.Sprintf("%s%s/reserve?num_tickets=%d", c.reserveURL, eventID, numTickets), "", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response struct {
		Tickets struct {
			EventID   string   `json:"event_id"`
			TicketIDs []string `json:"ticket_ids"`
		} `json:"tickets"`
	}

	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return response.Tickets.TicketIDs, nil
}

func simulateUserRequests(client *Client) {
    var wg sync.WaitGroup

    for i := 0; i < 5; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            events := client.GetEvents()
            fmt.Println("Events:", events)
        }()
    }

    wg.Wait()
}

func readUserInput(commandChan chan string) {
    scanner := bufio.NewScanner(os.Stdin)

    for scanner.Scan() {
        commandChan <- scanner.Text()
    }
    close(commandChan)
}

func processCommands(client *Client, commandChan chan string) {
    for command := range commandChan {
        if command == "getEvents" {
            go func() {
                client.GetEvents()
            }()
        } else if strings.HasPrefix(command, "bookEvents") {
            parts := strings.Split(command, " ")
            if len(parts) != 3 {
                fmt.Println("Invalid command format. Usage: bookEvents <eventID> <numTickets>")
                continue
            }
            eventID := parts[1]
            numTicketsStr := parts[2]
            numTickets, err := strconv.Atoi(numTicketsStr)
            if err != nil {
                fmt.Println("Invalid number of tickets:", err)
                continue
            }
            ticketIDs, err := client.BookTickets(eventID, numTickets)
            if err != nil {
                fmt.Println("Error reserving tickets")
            } else {
                fmt.Println("Reserved ticket IDs:", len(ticketIDs))
            }
        } else {
            fmt.Println("Invalid command")
        }
    }
}


func main() {
    fmt.Println("Client is running")

    client := createClient()

    commandChan := make(chan string)

    go readUserInput(commandChan)

    processCommands(client, commandChan)

    fmt.Println("Exiting...")
}