package server

import (
	"fmt"
	"log"
	"net/http"
	"ticket_reservation/src/TicketService"
)

func startServer(port string) {
	log.Println("Server is running on port", port)
	fmt.Println("ُServer Started!\nServer is running on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func loadData(ts *TicketService.TicketService) {
	ts.LoadEvents()
	ts.LoadTickets()
	log.Println("Data loaded.")
}

func Run(port string) {
	const maxClients = 10
	sema := make(chan struct{}, maxClients)
	createServerLogFile()

	ts := TicketService.TicketService{}

	loadData(&ts)

	http.HandleFunc("/", GetHomePageHandler)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		GetListEventsHandler(w, r, &ts)
	})

	http.HandleFunc("/events/{eventID}/reserve", func(w http.ResponseWriter, r *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		BookTicketsHandler(w, r, &ts)
	})

	http.HandleFunc("/events/create", func(w http.ResponseWriter, r *http.Request) {
		sema <- struct{}{}
		defer func() { <-sema }()
		CreateEventHandler(w, r, &ts)
	})

	startServer(port)
}
