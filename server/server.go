package server

import (
	"time"
	"fmt"
	"sync"
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
	fmt.Println("Server is running on port 8080")
}
