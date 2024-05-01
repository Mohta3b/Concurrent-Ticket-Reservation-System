package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"ticket_reservation/src/TicketService"
	"time"
)

func GetListEventsHandler(w http.ResponseWriter, r *http.Request, ts *TicketService.TicketService) {
	log.Println("GET /events")
	events := ts.ListEvents()
	fmt.Fprintf(w, "{\"events\": [")
	for i, event := range events {
		fmt.Fprintf(w, "{\"id\": \"%s\", \"name\": \"%s\", \"date\": \"%s\", \"total_tickets\": %d, \"available_tickets\": %d}", event.ID, event.Name, event.Date.Format("2006-01-02"), event.TotalTickets, event.AvailableTickets)
		if i < len(events)-1 {
			fmt.Fprintf(w, ", ")
		}
	}
	fmt.Fprintf(w, "]}")
}

func GetHomePageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the ticket service")
}

func BookTicketsHandler(w http.ResponseWriter, r *http.Request, ts *TicketService.TicketService) {
	type TicketResponse struct {
		Tickets struct {
			EventID   string   `json:"event_id"`
			TicketIDs []string `json:"ticket_ids"`
		} `json:"tickets"`
	}

	parsedURL, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing URL: %v", err), http.StatusBadRequest)
		return
	}

	pathParts := strings.Split(parsedURL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL format: eventID not found", http.StatusBadRequest)
		return
	}
	eventID := pathParts[2]

	numTicketsStr := r.URL.Query().Get("num_tickets")
	numTickets, err := strconv.Atoi(numTicketsStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing num_tickets: %v", err), http.StatusBadRequest)
		return
	}

	ticketIDs, err := ts.BookTickets(eventID, numTickets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error booking tickets: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(ticketIDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
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
		return
	}

	name := req.Name
	dateStr := req.Date
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing date: %v", err), http.StatusBadRequest)
		return
	}
	totalTickets, err := strconv.Atoi(req.Total)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing total_tickets: %v", err), http.StatusBadRequest)
		return
	}

	event, err := ts.CreateEvent(name, date, totalTickets)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating event: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(event)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
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
