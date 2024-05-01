package server

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"ticket_reservation/src/TicketService"
	"time"
)

func GetHomePageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("USER: GET /")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the ticket reservation system\n"))
}

func GetListEventsHandler(w http.ResponseWriter, r *http.Request, ts *TicketService.TicketService) {
	// define response body
	log.Println("USER: GET /events")
	events := ts.ListEvents()
	
	// save events in json response w.Write it
	response, err := json.Marshal(events)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		log.Println("Error encoding response: %v", err)
		return
	}

	// write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	// print response body
	// log.Println("List of events:", string(response))
}

func BookTicketsHandler(w http.ResponseWriter, r *http.Request, ts *TicketService.TicketService) {
	
	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusBadRequest)
		log.Println("Error parsing URL: %v", err)
		return
	}

	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format: eventID not found", http.StatusBadRequest)
		log.Println("Invalid URL format: eventID not found")
		return
	}
	eventID := pathParts[2]

	numTicketsStr := r.URL.Query().Get("num_tickets")
	numTickets, err := strconv.Atoi(numTicketsStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing num_tickets: %v", err), http.StatusBadRequest)
		log.Println("Error parsing num_tickets: %v", err)
		return
	}

	log.Println("USER: GET /events/" + eventID + "/tickets?num_tickets=" + numTicketsStr)

	ticketIDs, err := ts.BookTickets(eventID, numTickets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error booking tickets: %v", err), http.StatusInternalServerError)
		log.Println("Error booking tickets: %v", err)
		return
	}

	response, err := json.Marshal(ticketIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		log.Println("Error encoding response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

	log.Printf("Booked %v tickets for event %s", ticketIDs, eventID)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request, ts *TicketService.TicketService) {
	type CreateEventResponse struct {
		Name  string `json:"name"`
		Date  string `json:"date"`
		Total string `json:"total_tickets"`
	}

	var req CreateEventResponse
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
		log.Println("Error decoding request: %v", err)
		return
	}

	name := req.Name
	dateStr := req.Date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing date: %v", err), http.StatusBadRequest)
		log.Println("Error parsing date: %v", err)
		return
	}
	totalTickets, err := strconv.Atoi(req.Total)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing total_tickets: %v", err), http.StatusBadRequest)
		log.Println("Error parsing total_tickets: %v", err)
		return
	}

	event, err := ts.CreateEvent(name, date, totalTickets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating event: %v", err), http.StatusInternalServerError)
		log.Println("Error creating event: %v", err)
		return
	}

	log.Printf("USER: POST /events\nCreated event %s with %d tickets", event.ID, event.TotalTickets)

	response, err := json.Marshal(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		log.Println("Error encoding response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Created event %s with %d tickets", event.ID, event.TotalTickets)
}


func createServerLogFile() {
	// create server_log.txt file in ./data directory
	logFile, err := os.Create("./data/server_log.txt")

	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)

	log.Println("!!! Server Log !!!")
}