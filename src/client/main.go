package main

// import (
// 	"ticket_reservation/src/client"
// )

const PORT = ":5050"

func main() {
    serverURL := "http://localhost" + PORT
    eventURL := "http://localhost" + PORT + "/events"
    reserveURL := "http://localhost" + PORT + "/events/"
	// eventURL := "http://localhost:5050/events"
    // bookURL := "http://localhost:5050/events/book"
	// createURL := "http://localhost:5050/events/create"
    client := NewClient(serverURL, eventURL, reserveURL)
    
    ConnectToServer(client)
}
