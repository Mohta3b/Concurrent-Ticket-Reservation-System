package server

import (
	"fmt"
	"log"
	"net/http"
	"ticket_reservation/src/TicketService"
	// "ticket_reservation/src/Event"
	// "ticket_reservation/src/Ticket"
)

func startServer() {
	log.Println("Server is running on port 5050")
	log.Fatal(http.ListenAndServe(":5050", nil))
}


func Run() {
	fmt.Println("Server is running on port 8080")

	ts := TicketService.TicketService{}
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		GetListEventsHandler(w, r, &ts)
	})

	http.HandleFunc("/", GetHomePageHandler)

	http.HandleFunc("/events/{eventID}/reserve", func(w http.ResponseWriter, r *http.Request) {
		BookTicketsHandler(w, r, &ts)
	})
	
	http.HandleFunc("/events/create", func(w http.ResponseWriter, r *http.Request) {
		CreateEventHandler(w, r, &ts)
	})

	startServer()
}
