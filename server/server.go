package server

import (
	"time"
	"fmt"
	"sync"
	"net/http"
)

type Event struct {
	ID   				string
	Name 				string
	Date 				time.Time
	TotalTickets		int
	AvailableTickets	int
}

type Ticket struct {
	ID		string
	EventID	string
}

type TicketService struct {
	events sync.Map
	tickets sync.Map
}


// run server using http server
func Run() {
	ts := &TicketService{}
	// recieve requests
	// home page
	http.HandleFunc("/", getHomePageHandler)
	
	// list events: request to http://localhost:8080/events
	http.HandleFunc("/events", ts.getListEventsHandler)

	// reserve tickets (eventID and number of tickets are passed): request to http://localhost:8080/events/{eventID}/reserve?num_tickets={num_tickets}
	http.HandleFunc("/events/{eventID}/reserve", ts.reserveTicketsHandler)
	
	// start server
	http.ListenAndServe(":8080", nil)

	fmt.Println("Server is running on port 8080")
	
}
