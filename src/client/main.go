package main

// import (
// 	"ticket_reservation/src/client"
// )

func main() {
    eventURL := "http://localhost:5050/events"
    reserveURL := "http://localhost:5050/events/"
	// eventURL := "http://localhost:5050/events"
    // bookURL := "http://localhost:5050/events/book"
	// createURL := "http://localhost:5050/events/create"
    client := NewClient(eventURL, reserveURL)
	GetInput(client)
}
