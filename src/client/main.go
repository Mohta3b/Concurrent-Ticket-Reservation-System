package main

// import (
// 	"ticket_reservation/src/client"
// )

func main() {
    eventURL := "http://localhost:5050/events"
    reserveURL := "http://localhost:5050/events/"
    client := NewClient(eventURL, reserveURL)
	GetInput(client)
}
